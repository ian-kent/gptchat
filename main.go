package main

import (
	"fmt"
	"github.com/ian-kent/gptchat/module"
	"github.com/ian-kent/gptchat/module/memory"
	"github.com/ian-kent/gptchat/module/plugin"
	"github.com/ian-kent/gptchat/ui"
	openai "github.com/sashabaranov/go-openai"
	"os"
	"strconv"
)

var client = openai.NewClient(os.Getenv("OPENAI_API_KEY"))

func init() {
	module.Load(client, []module.Module{
		&memory.Module{},
		&plugin.Module{},
	}...)
	if err := module.LoadCompiledPlugins(); err != nil {
		fmt.Printf("error loading compiled plugins: %s", err)
		os.Exit(1)
	}
}

func main() {
	debugMode := false
	debugEnv := os.Getenv("GPT_DEBUG")
	if debugEnv != "" {
		v, err := strconv.ParseBool(debugEnv)
		if err != nil {
			ui.Warn(fmt.Sprintf("error parsing GPT_DEBUG: %s", err.Error()))
		} else {
			debugMode = v
		}
	}

	ui.Welcome(
		`Welcome to the GPT-4 client.`,
		`You can talk directly to GPT-4, or you can use /commands to interact with the client.

Use /help to see a list of available commands.`)

	chatLoop(debugMode)
}
