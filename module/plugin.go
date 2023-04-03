package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ian-kent/gptchat/config"
	"github.com/ian-kent/gptchat/ui"
	"github.com/sashabaranov/go-openai"
	"os"
	"plugin"
	"strings"
)

type Plugin interface {
	ID() string
	Example() string
	Execute(map[string]any) (map[string]any, error)
}

type pluginLoader struct {
	plugin Plugin
}

func (p pluginLoader) Load(config.Config, *openai.Client) error {
	return nil
}
func (p pluginLoader) UpdateConfig(config.Config) {}
func (p pluginLoader) ID() string {
	return p.plugin.ID()
}
func (p pluginLoader) Prompt() string {
	return p.plugin.Example()
}
func (p pluginLoader) Execute(args, body string) (string, error) {
	input := make(map[string]any)
	if body != "" {
		err := json.Unmarshal([]byte(body), &input)
		if err != nil {
			return "", fmt.Errorf("plugin body must be valid json: %s", err)
		}
	}

	output, err := p.plugin.Execute(input)
	if err != nil {
		return "", fmt.Errorf("error executing plugin: %s", err)
	}

	b, err := json.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("error converting plugin output to json: %s", err)
	}

	return string(b), nil
}

func GetModuleForPlugin(p Plugin) Module {
	return pluginLoader{p}
}

func LoadCompiledPlugins() error {
	pluginPath := "./module/plugin/compiled/"
	entries, err := os.ReadDir(pluginPath)
	if err != nil {
		return fmt.Errorf("error loading compiled plugins: %s", err)
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".so") {
			continue
		}

		loadedPlugin, err := OpenPlugin(pluginPath + entry.Name())
		if err != nil {
			ui.Warn(fmt.Sprintf("error opening plugin: %s", err))
			continue
		}

		pluginID := loadedPlugin.ID()
		if IsLoaded(pluginID) {
			ui.Warn(fmt.Sprintf("plugin with this ID is already loaded: %s", err))
			continue
		}

		err = LoadPlugin(GetModuleForPlugin(loadedPlugin))
		if err != nil {
			ui.Warn(fmt.Sprintf("error loading plugin: %s", err))
			continue
		}
	}

	return nil
}

func OpenPlugin(path string) (Plugin, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error loading plugin: %s", err)
	}

	apiSymbol, err := p.Lookup("Plugin")
	if err != nil {
		return nil, fmt.Errorf("error finding plugin implementation: %s", err)
	}

	// Cast the symbol to the ScriptAPI interface
	api, ok := apiSymbol.(*Plugin)
	if !ok {
		return nil, errors.New("plugin does not implement the Plugin interface")
	}

	loadedPlugin := *api
	return loadedPlugin, nil
}
