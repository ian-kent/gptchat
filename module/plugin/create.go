package plugin

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/ian-kent/gptchat/config"
	"github.com/ian-kent/gptchat/module"
	"github.com/ian-kent/gptchat/ui"
	"github.com/ian-kent/gptchat/util"
	openai "github.com/sashabaranov/go-openai"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Module struct {
	cfg    config.Config
	client *openai.Client
}

func (m *Module) Load(cfg config.Config, client *openai.Client) error {
	m.cfg = cfg
	m.client = client
	return nil
}

func (m *Module) UpdateConfig(cfg config.Config) {
	m.cfg = cfg
}

func (m *Module) Prompt() string {
	return newPluginPrompt
}

func (m *Module) ID() string {
	return "plugin"
}

func (m *Module) Execute(args, body string) (string, error) {
	parts := strings.SplitN(args, " ", 2)
	cmd := parts[0]
	if len(parts) > 1 {
		args = parts[1]
	}

	switch cmd {
	case "create":
		return m.createPlugin(args, body)
	default:
		return "", errors.New(fmt.Sprintf("%s not implemented", args))
	}
}

func (m *Module) createPlugin(id, body string) (string, error) {
	body = strings.TrimSpace(body)
	if len(body) == 0 {
		return "", errors.New("plugin source not found")
	}

	if !strings.HasPrefix(body, "{") || !strings.HasSuffix(body, "}") {
		return "", errors.New("plugin source must be between {} in '/plugin create plugin-id {}' command")
	}

	id = strings.TrimSpace(id)
	if id == "" {
		return "", errors.New("plugin id is invalid")
	}

	if module.IsLoaded(id) {
		return "", errors.New("a plugin with this id already exists")
	}

	source := strings.TrimPrefix(strings.TrimSuffix(body, "}"), "{")

	err := os.Mkdir("./module/plugin/source/"+id, 0777)
	if err != nil {
		return "", fmt.Errorf("error creating directory: %s", err)
	}

	sourcePath := "./module/plugin/source/" + id + "/plugin.go"
	err = ioutil.WriteFile("./module/plugin/source/"+id+"/plugin.go", []byte(source), 0644)
	if err != nil {
		return "", fmt.Errorf("error writing source file: %s", err)
	}

	if m.cfg.IsSupervisedMode() {
		fmt.Println("============================================================")
		fmt.Println()
		ui.Warn("⚠️ GPT written plugins are untrusted code from the internet")
		fmt.Println()
		fmt.Println("You should review this code before allowing it to be compiled and executed.")
		fmt.Println()
		fmt.Println("If you allow this action, GPT is able to execute code with the same permissions as your user.")
		fmt.Println()
		color.New(color.FgHiWhite, color.Bold).Println("This is potentially dangerous.")
		fmt.Println()
		fmt.Println("The source code GPT has written can be found here:")
		fmt.Println(sourcePath)
		fmt.Println()
		confirmation := ui.PromptInput("Enter 'confirm' to confirm, anything else will block:")
		if confirmation != "confirm" {
			fmt.Println()
			fmt.Println("============================================================")
			return "The user has prevented you from running this code", errors.New(confirmation)
		}
		fmt.Println()
		fmt.Println("============================================================")
		fmt.Println()
	}

	pluginPath := "./module/plugin/compiled/" + id + ".so"
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", pluginPath, sourcePath)
	if b, err := cmd.CombinedOutput(); err != nil {
		return string(b), fmt.Errorf("error compiling plugin: %s", err)
	}

	loadedPlugin, err := module.OpenPlugin(pluginPath)
	if err != nil {
		return "", fmt.Errorf("error opening plugin: %s", err)
	}

	// Call the functions provided by the plugin
	compiledID := loadedPlugin.ID()
	if id != compiledID {
		return "", errors.New("ID() does not return the ID specified in the '/plugin create <plugin-id>' command")
	}

	err = module.LoadPlugin(module.GetModuleForPlugin(loadedPlugin))
	if err != nil {
		return "", fmt.Errorf("error loading plugin: %s", err)
	}

	return `Great! Your plugin has been loaded successfully.

Why don't you check the /help command to see if your new plugin is available.'`, nil
}

var newPluginPrompt = `You can add new plugins which you can call using a slash command.

They're written in Go, so all you need to do is create a new struct which implements the correct interface.

The interface you need to implement is:

` + util.TripleQuote + `
type Plugin interface {
	Example() string
	Execute(input map[string]any) (map[string]any, error)
}
` + util.TripleQuote + `

You don't need to write any supporting code like the main function, you only need to implement the struct.

Here's the full code for the "add 1" plugin you can use to guide your output:
` + util.TripleQuote + `
package main

import "github.com/ian-kent/gptchat/module"

var Plugin module.Plugin = AddOne{}

type AddOne struct{}

func (c AddOne) ID() string {
	return "add-one"
}

func (c AddOne) Example() string {
	return ` + util.SingleQuote + `/add-one {
	"value": 5
}` + util.SingleQuote + `
}

func (c AddOne) Execute(input map[string]any) (map[string]any, error) {
	value, ok := input["value"].(int)
	if !ok {
		return nil, nil
	}

	value = value + 1

	return map[string]any{
		"result": value,
	}, nil
}
` + util.TripleQuote + `

It's best if the plugins you create don't have any external dependencies. You can call external APIs if you want to, but you should avoid APIs which require authentication since you won't have the required access.

Your plugin must import the module package and must define a package variable named 'Plugin', just like with the AddOne example. The result of the Execute function you implement must return either a value or an error.

The input to Execute is a map[string]any which you should assume is unmarshaled from JSON. This means you must use appropriate data types, for example a float64 when working with numbers.

To create a plugin, you should use the "/plugin create <plugin-id> {}" command, for example:

` + util.TripleQuote + `
/plugin create add-one {
	package main

	// the rest of your plugin source here
}
` + util.TripleQuote + `

Your code inside the '/plugin create' body must be valid Go code which can compile without any errors. Do not include quotes or attempt to use a JSON body.`
