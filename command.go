package main

import (
	"fmt"
	"github.com/ian-kent/gptchat/ui"
	"os"
	"strings"
)

type slashCommandResult struct {
	// the prompt to send to the AI
	prompt string

	// retry tells the client to resend the most recent conversation
	retry bool

	// resetConversation will reset the conversation to its original state,
	// forgetting the conversation history
	resetConversation bool

	// toggleDebugMode will switch between debug on and debug off
	toggleDebugMode bool

	// toggleSupervisedMode will switch between supervised mode on and off
	toggleSupervisedMode bool
}

type slashCommand struct {
	command string
	fn      func(string) (bool, *slashCommandResult)
}

var slashCommands = []slashCommand{
	{
		command: "exit",
		fn: func(s string) (bool, *slashCommandResult) {
			os.Exit(0)
			return true, nil
		},
	},
	{
		command: "retry",
		fn: func(s string) (bool, *slashCommandResult) {
			return true, &slashCommandResult{
				retry: true,
			}
		},
	},
	{
		command: "reset",
		fn: func(s string) (bool, *slashCommandResult) {
			return true, &slashCommandResult{
				resetConversation: true,
			}
		},
	},
	{
		command: "debug",
		fn: func(s string) (bool, *slashCommandResult) {
			return true, &slashCommandResult{
				toggleDebugMode: true,
			}
		},
	},
	{
		command: "supervisor",
		fn: func(s string) (bool, *slashCommandResult) {
			return true, &slashCommandResult{
				toggleSupervisedMode: true,
			}
		},
	},
	{
		command: "example",
		fn:      exampleCommand,
	},
}

func helpCommand(string) (bool, *slashCommandResult) {
	result := "The following commands are available:\n"
	for _, e := range slashCommands {
		result += fmt.Sprintf("\n    /%s", e.command)
	}

	ui.PrintChat(ui.App, result)

	return true, nil
}

func parseSlashCommand(input string) (ok bool, result *slashCommandResult) {
	if !strings.HasPrefix(input, "/") {
		return false, nil
	}

	input = strings.TrimPrefix(input, "/")

	if input == "help" {
		return helpCommand(input)
	}

	parts := strings.SplitN(input, " ", 2)
	var cmd, args string
	cmd = parts[0]
	if len(parts) > 1 {
		args = parts[1]
	}

	for _, command := range slashCommands {
		if command.command == cmd {
			return command.fn(args)
		}
	}

	return false, nil
}

type example struct {
	id, prompt string
}

var examples = []example{
	{
		id:     "1",
		prompt: "I want you to generate 5 random numbers and add them together.",
	},
	{
		id:     "2",
		prompt: "I want you to generate 5 random numbers. Multiply the first and second number, then add the result to the remaining numbers.",
	},
	{
		id:     "3",
		prompt: "I want you to generate 2 random numbers. Add them together then multiply the result by -1.",
	},
	{
		id:     "4",
		prompt: "Can you summarise the tools you have available?",
	},
	{
		id:     "5",
		prompt: "Can you suggest a task which might somehow use all of the available tools?",
	},
}

func exampleCommand(args string) (bool, *slashCommandResult) {
	for _, e := range examples {
		if e.id == args {
			return true, &slashCommandResult{
				prompt: e.prompt,
			}
		}
	}

	result := "The following examples are available:"
	for _, e := range examples {
		result += fmt.Sprintf("\n\n/example %s\n        %s", e.id, e.prompt)
	}

	ui.PrintChat(ui.App, result)

	return true, nil
}
