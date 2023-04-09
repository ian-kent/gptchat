package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ian-kent/gptchat/config"
	"github.com/ian-kent/gptchat/module"
	"github.com/ian-kent/gptchat/module/memory"
	"github.com/ian-kent/gptchat/module/plugin"
	"github.com/ian-kent/gptchat/ui"
	openai "github.com/sashabaranov/go-openai"
)

var client *openai.Client
var cfg = config.New()

func init() {
	openaiAPIKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if openaiAPIKey == "" {
		ui.Warn("You haven't configured an OpenAI API key")
		fmt.Println()
		if !ui.PromptConfirm("Do you have an API key?") {
			ui.Warn("You'll need an API key to use GPTChat")
			fmt.Println()
			fmt.Println("* You can get an API key at https://platform.openai.com/account/api-keys")
			fmt.Println("* You can get join the GPT-4 API waitlist at https://openai.com/waitlist/gpt-4-api")
			os.Exit(1)
		}

		openaiAPIKey = ui.PromptInput("Enter your API key:")
		if openaiAPIKey == "" {
			fmt.Println("")
			ui.Warn("You didn't enter an API key.")
			os.Exit(1)
		}
	}

	cfg = cfg.WithOpenAIAPIKey(openaiAPIKey)

	openaiAPIModel := strings.TrimSpace(os.Getenv("OPENAI_API_MODEL"))

	if openaiAPIModel == "" {
		ui.Warn("You haven't configured an OpenAI API model, defaulting to GPT4")

		openaiAPIModel = openai.GPT4
	}

	cfg = cfg.WithOpenAIAPIModel(openaiAPIModel)

	supervisorMode := os.Getenv("GPTCHAT_SUPERVISOR")
	switch strings.ToLower(supervisorMode) {
	case "disabled":
		ui.Warn("Supervisor mode is disabled")
		cfg = cfg.WithSupervisedMode(false)
	default:
	}

	debugEnv := os.Getenv("GPTCHAT_DEBUG")
	if debugEnv != "" {
		v, err := strconv.ParseBool(debugEnv)
		if err != nil {
			ui.Warn(fmt.Sprintf("error parsing GPT_DEBUG: %s", err.Error()))
		} else {
			cfg = cfg.WithDebugMode(v)
		}
	}

	client = openai.NewClient(openaiAPIKey)

	module.Load(cfg, client, []module.Module{
		&memory.Module{},
		&plugin.Module{},
	}...)

	if err := module.LoadCompiledPlugins(); err != nil {
		ui.Warn(fmt.Sprintf("error loading compiled plugins: %s", err))
	}
}

func main() {
	ui.Welcome(
		`Welcome to the GPT client.`,
		`You can talk directly to GPT, or you can use /commands to interact with the client.

Use /help to see a list of available commands.`)

	chatLoop(cfg)
}
