package config

import (
	"errors"
	"slices"

	"github.com/Mirai3103/Project-Re-ENE/config/tts"
)

var supportedTTSProviders = []string{"elevenlabs"}

type TTSConfig struct {
	Provider         string                `yaml:"provider"`
	ElevenLabsConfig *tts.ElevenLabsConfig `yaml:"eleven_labs_config"`
}

func (c *TTSConfig) Validate() error {
	if !slices.Contains(supportedTTSProviders, c.Provider) {
		return errors.New("tts provider is not supported: " + c.Provider)
	}
	switch c.Provider {
	case "elevenlabs":
		return c.ElevenLabsConfig.Validate()
	default:
		return errors.New("tts provider is not supported: " + c.Provider)
	}
}

func getDefaultTTSConfig() *TTSConfig {
	return &TTSConfig{
		Provider:         "elevenlabs",
		ElevenLabsConfig: tts.GetDefaultElevenLabsConfig(),
	}
}
