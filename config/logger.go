package config

import (
	"errors"
	"slices"
)

var supportedLoggerModes = []string{"console", "file"}

type LoggerConfig struct {
	Mode     string `yaml:"mode"`
	Level    string `yaml:"level"`
	FilePath string `yaml:"file_path"`
}

func (c *LoggerConfig) Validate() error {
	if c.Mode == "" {
		return errors.New("mode is required")
	}
	if c.Level == "" {
		return errors.New("level is required")
	}
	if !slices.Contains(supportedLoggerModes, c.Mode) {
		return errors.New("mode is not supported")
	}
	if c.FilePath == "" && c.Mode == "file" {
		return errors.New("file_path is required")
	}
	return nil
}

func getDefaultLoggerConfig() *LoggerConfig {

	return &LoggerConfig{
		Mode:     "console",
		Level:    "info",
		FilePath: "logs/app.log",
	}
}
