package main

import (
	"context"
	"fmt"
	"github.com/ian-kent/gptchat/config"
	"github.com/ian-kent/gptchat/module"
	"github.com/ian-kent/gptchat/parser"
	"github.com/ian-kent/gptchat/ui"
	"github.com/ian-kent/gptchat/util"
	"github.com/sashabaranov/go-openai"
	"strings"
	"time"
)

func chatLoop(cfg config.Config) {
RESET:
	appendMessage(openai.ChatMessageRoleSystem, systemPrompt)
	if cfg.IsDebugMode() {
		ui.PrintChatDebug(ui.System, systemPrompt)
	}

	var skipUserInput = true
	appendMessage(openai.ChatMessageRoleUser, openingPrompt)
	if cfg.IsDebugMode() {
		ui.PrintChatDebug(ui.User, openingPrompt)
	}

	if !cfg.IsDebugMode() {
		ui.PrintChat(ui.App, "Setting up the chat environment, please wait for GPT to respond - this may take a few moments.")
	}

	var i int
	for {
		i++

		if !skipUserInput {
			input := ui.PromptChatInput()
			var echo bool

			ok, result := parseSlashCommand(input)
			if ok {
				// the command was handled but returned nothing
				// to send to the AI, let's prompt the user again
				if result == nil {
					continue
				}

				if result.resetConversation {
					resetConversation()
					goto RESET
				}

				// if the result is a retry, we can just send the
				// same request to GPT again
				if result.retry {
					skipUserInput = true
					goto RETRY
				}

				if result.toggleDebugMode {
					cfg = cfg.WithDebugMode(!cfg.IsDebugMode())
					module.UpdateConfig(cfg)
					if cfg.IsDebugMode() {
						ui.PrintChat(ui.App, "Debug mode is now enabled")
					} else {
						ui.PrintChat(ui.App, "Debug mode is now disabled")
					}
					continue
				}

				if result.toggleSupervisedMode {
					cfg = cfg.WithSupervisedMode(!cfg.IsSupervisedMode())
					module.UpdateConfig(cfg)
					if cfg.IsSupervisedMode() {
						ui.PrintChat(ui.App, "Supervised mode is now enabled")
					} else {
						ui.PrintChat(ui.App, "Supervised mode is now disabled")
					}
					continue
				}

				// we have a prompt to give to the AI, let's do that
				if result.prompt != "" {
					input = result.prompt
					echo = true
				}
			}

			if echo {
				ui.PrintChat(ui.User, input)
				echo = false
			}

			appendMessage(openai.ChatMessageRoleUser, input)
		}

		skipUserInput = false

	RETRY:

		// Occasionally include the interval prompt
		if i%5 == 0 {
			interval := intervalPrompt()
			appendMessage(openai.ChatMessageRoleSystem, interval)
			if cfg.IsDebugMode() {
				ui.PrintChatDebug(ui.System, interval)
			}
		}

		attempts := 1
	RATELIMIT_RETRY:
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    openai.GPT4,
				Messages: conversation,
			},
		)
		if err != nil {
			if strings.HasPrefix(err.Error(), "error, status code: 429") && attempts < 5 {
				attempts++
				ui.Error("rate limited, trying again in 1 second", err)
				time.Sleep(time.Second)
				goto RATELIMIT_RETRY
			}

			ui.Error("ChatCompletion failed", err)
			if ui.PromptConfirm("Would you like to try again?") {
				goto RATELIMIT_RETRY
			}

			continue
		}

		response := resp.Choices[0].Message.Content
		appendMessage(openai.ChatMessageRoleAssistant, response)
		if cfg.IsDebugMode() {
			ui.PrintChat(ui.AI, response)
		}

		parseResult := parser.Parse(response)

		if !cfg.IsDebugMode() && parseResult.Chat != "" {
			ui.PrintChat(ui.AI, parseResult.Chat)
		}

		for _, command := range parseResult.Commands {
			ok, result := module.ExecuteCommand(command.Command, command.Args, command.Body)
			if ok {
				// we had at least one AI command so we're going to respond automatically,
				// no need to ask for user input
				skipUserInput = true

				if result.Error != nil {
					msg := fmt.Sprintf(`An error occurred executing your command.

The command was:
`+util.TripleQuote+`
%s
`+util.TripleQuote+`

The error was:
`+util.TripleQuote+`
%s
`+util.TripleQuote, command.String(), result.Error.Error())

					if result.Prompt != "" {
						msg += fmt.Sprintf(`

The command provided this additional output:
`+util.TripleQuote+`
%s
`+util.TripleQuote, result.Prompt)
					}

					appendMessage(openai.ChatMessageRoleSystem, msg)
					if cfg.IsDebugMode() {
						ui.PrintChatDebug(ui.Module, msg)
					}
					continue
				}

				commandResult := fmt.Sprintf(`Your command returned some output.

The command was:
`+util.TripleQuote+`
%s
`+util.TripleQuote+`

The output was:

%s`, command.String(), result.Prompt)
				appendMessage(openai.ChatMessageRoleSystem, commandResult)

				if cfg.IsDebugMode() {
					ui.PrintChatDebug(ui.Module, commandResult)
				}
				continue
			}
		}
	}
}
