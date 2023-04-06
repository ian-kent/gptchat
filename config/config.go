package config

type Config struct {
	openaiAPIKey   string
	openaiAPIModel string

	supervisedMode bool
	debugMode      bool
}

func New() Config {
	return Config{
		openaiAPIKey:   "",
		openaiAPIModel: "",
		supervisedMode: true,
		debugMode:      false,
	}
}

func (c Config) OpenAIAPIModel() string {
	return c.openaiAPIModel
}

func (c Config) OpenAIAPIKey() string {
	return c.openaiAPIKey
}

func (c Config) IsSupervisedMode() bool {
	return c.supervisedMode
}

func (c Config) IsDebugMode() bool {
	return c.debugMode
}

func (c Config) WithOpenAIAPIKey(apiKey string) Config {
	c.openaiAPIKey = apiKey
	return c
}

func (c Config) WithSupervisedMode(supervisedMode bool) Config {
	c.supervisedMode = supervisedMode
	return c
}

func (c Config) WithDebugMode(debugMode bool) Config {
	c.debugMode = debugMode
	return c
}

func (c Config) WithOpenAIAPIModel(apiModel string) Config {
	c.openaiAPIModel = apiModel
	return c
}
