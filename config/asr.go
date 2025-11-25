package config

import (
	"errors"
	"slices"

	"github.com/Mirai3103/Project-Re-ENE/config/asr"
)

var supportedASRProviders = []string{"elevenlabs"}

type ASRConfig struct {
	Provider         string                `yaml:"provider"`
	ElevenLabsConfig *asr.ElevenLabsConfig `yaml:"eleven_labs_config"`
	InputDevice      string                `yaml:"input_device"`
}

func (c *ASRConfig) Validate() error {
	if !slices.Contains(supportedASRProviders, c.Provider) {
		return errors.New("asr provider is not supported: " + c.Provider)
	}
	switch c.Provider {
	case "elevenlabs":
		return c.ElevenLabsConfig.Validate()
	default:
		return errors.New("asr provider is not supported: " + c.Provider)
	}
}

func GetDefaultASRConfig() *ASRConfig {
	return &ASRConfig{
		Provider:         "elevenlabs",
		ElevenLabsConfig: asr.GetDefaultElevenLabsConfig(),
		InputDevice:      "default",
	}
}
