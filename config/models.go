package config

import (
	"errors"
	"os"
)

type ModelsConfig struct {
	ModelDir string `yaml:"model_dir"`
}

func (m *ModelsConfig) Validate() error {
	if m.ModelDir == "" {
		return errors.New("model_dir is required")
	}
	if err := os.MkdirAll(m.ModelDir, 0755); err != nil {
		return errors.New("failed to create model dir")
	}
	return nil
}

func getDefaultModelsConfig() *ModelsConfig {
	return &ModelsConfig{
		ModelDir: "./resources/live2d",
	}
}
