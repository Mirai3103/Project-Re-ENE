package llm

import "errors"

type GeminiConfig struct {
	BaseLLMConfig
	APIKey      string  `yaml:"api_key"`
	Model       string  `yaml:"model"`
	BaseURL     string  `yaml:"base_url"`
	Temperature float32 `yaml:"temperature"`
}

func (g *GeminiConfig) Validate() error {
	if g.APIKey == "" {
		return errors.New("api_key is required")
	}
	if g.Model == "" {
		return errors.New("model is required")
	}
	return nil
}

func GetDefaultGeminiConfig() *GeminiConfig {
	return &GeminiConfig{
		APIKey:      "",
		Model:       "gemini-2.0-flash",
		Temperature: 0.7,
	}
}
