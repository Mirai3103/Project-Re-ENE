package tts

import (
	"errors"
)

type ElevenLabsConfig struct {
	BaseTTSConfig
	APIKey  string `yaml:"api_key"`
	ModelID string `yaml:"model_id"`
	VoiceID string `yaml:"voice_id"`
}

func (e *ElevenLabsConfig) Validate() error {
	if e.APIKey == "" {
		return errors.New("api_key is required")
	}
	if e.ModelID == "" {
		return errors.New("model_id is required")
	}
	if e.VoiceID == "" {
		return errors.New("voice_id is required")
	}
	return nil
}

func GetDefaultElevenLabsConfig() *ElevenLabsConfig {
	return &ElevenLabsConfig{
		APIKey:  "",
		ModelID: "",
		VoiceID: "21m00Tcm4TlvDq8ikWAM",
	}
}
