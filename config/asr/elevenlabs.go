package asr

import (
	"errors"
)

type ElevenLabsConfig struct {
	Name         string `yaml:"name"`
	APIKey       string `yaml:"api_key"`
	ModelID      string `yaml:"model_id"`
	LanguageCode string `yaml:"language_code"`
}

func (e *ElevenLabsConfig) Validate() error {
	if e.APIKey == "" {
		return errors.New("api_key is required")
	}
	if e.ModelID == "" {
		return errors.New("model_id is required")
	}
	return nil
}

func GetDefaultElevenLabsConfig() *ElevenLabsConfig {
	return &ElevenLabsConfig{
		APIKey:       "",
		ModelID:      "scribe_v1",
		LanguageCode: "auto",
	}
}
