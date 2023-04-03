package module

import (
	"errors"
	"fmt"
	"github.com/ian-kent/gptchat/config"
	openai "github.com/sashabaranov/go-openai"
	"strings"
)

type Module interface {
	Load(config.Config, *openai.Client) error
	UpdateConfig(config.Config)
	ID() string
	Prompt() string
	Execute(args, body string) (string, error)
}

// IntervalPrompt allows a module to inject a prompt into the interval prompt
type IntervalPrompt interface {
	IntervalPrompt() string
}

var loadedModules = make(map[string]Module)

func Load(cfg config.Config, client *openai.Client, modules ...Module) error {
	for _, module := range modules {
		if err := module.Load(cfg, client); err != nil {
			return err
		}
		loadedModules[module.ID()] = module
	}
	return nil
}

func UpdateConfig(cfg config.Config) {
	for _, module := range loadedModules {
		_, ok := module.(pluginLoader)
		if ok {
			// GPT written plugins shouldn't have config, nothing to do
			continue
		}

		module.UpdateConfig(cfg)
	}
}

func IsLoaded(id string) bool {
	_, ok := loadedModules[id]
	return ok
}

func LoadPlugin(m Module) error {
	// a plugin doesn't have access to the openai client so it's safe to pass in nil here
	//
	// we also don't pass in the config since it may contain sensitive information that
	// we don't want GPT to have access to
	return Load(config.Config{}, nil, m)
}

type CommandResult struct {
	Error  error
	Prompt string
}

func HelpCommand() (bool, *CommandResult) {
	result := "Here are the commands you have available:\n\n"
	for _, mod := range loadedModules {
		result += fmt.Sprintf("    * /%s\n", mod.ID())
	}
	result += `
You can call commands using the /command syntax.

Calling a command without any additional arguments will explain it's usage. You should do this to learn how the command works.`

	return true, &CommandResult{
		Prompt: result,
	}
}

func ExecuteCommand(command, args, body string) (bool, *CommandResult) {
	if command == "/help" {
		return HelpCommand()
	}

	cmd := strings.TrimPrefix(command, "/")
	mod, ok := loadedModules[cmd]
	if !ok {
		return true, &CommandResult{
			Error: errors.New(fmt.Sprintf("Unrecognised command: %s", command)),
		}
	}

	if args == "" && body == "" {
		return true, &CommandResult{
			Prompt: mod.Prompt(),
		}
	}

	res, err := mod.Execute(args, body)
	if err != nil {
		return true, &CommandResult{
			Error: err,
		}
	}

	return true, &CommandResult{
		Prompt: res,
	}
}
