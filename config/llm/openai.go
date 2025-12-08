package llm

type OpenAIConfig struct {
	APIKey      string  `yaml:"api_key"`
	Model       string  `yaml:"model"`
	BaseURL     string  `yaml:"base_url"`
	Temperature float32 `yaml:"temperature"`
}

func GetDefaultOpenAIConfig() *OpenAIConfig {
	return &OpenAIConfig{
		APIKey:      "",
		Model:       "gpt-4o-mini",
		Temperature: 0.7,
		BaseURL:     "https://api.openai.com/v1",
	}
}
