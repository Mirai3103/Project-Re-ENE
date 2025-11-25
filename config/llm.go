package config

import (
	"errors"
	"slices"

	"github.com/Mirai3103/Project-Re-ENE/config/llm"
)

var supportedLLMProviders = []string{"gemini"}

type LLMConfig struct {
	Provider          string            `yaml:"provider"`
	IsAudioSupported  bool              `yaml:"is_audio_supported"`
	IsVisionSupported bool              `yaml:"is_vision_supported"`
	GeminiConfig      *llm.GeminiConfig `yaml:"gemini_config"`
	// todo: add other LLM configs
}

func (c *LLMConfig) Validate() error {
	if !slices.Contains(supportedLLMProviders, c.Provider) {
		return errors.New("llm provider is not supported: " + c.Provider)
	}
	switch c.Provider {
	case "gemini":
		return c.GeminiConfig.Validate()
	default:
		return errors.New("llm provider is not supported")
	}
}

func getDefaultLLMConfig() *LLMConfig {
	return &LLMConfig{
		Provider:          "gemini",
		IsAudioSupported:  true,
		IsVisionSupported: true,
		GeminiConfig:      llm.GetDefaultGeminiConfig(),
	}
}
