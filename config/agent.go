package config

import (
	"errors"

	tool "github.com/Mirai3103/Project-Re-ENE/config/tool"
)

type ShortTermMemoryConfig struct {
	MaxWindowSize int `yaml:"max_window_size"`
}

func (c *ShortTermMemoryConfig) Validate() error {
	if c.MaxWindowSize <= 0 {
		return errors.New("max_window_size must be greater than 0")
	}

	return nil
}

type AgentConfig struct {
	ShortTermMemoryConfig ShortTermMemoryConfig `yaml:"short_term_memory_config"`
	ToolsConfig           tool.ToolConfig       `yaml:"tools_config"`
}

func getDefaultAgentConfig() *AgentConfig {
	return &AgentConfig{
		ShortTermMemoryConfig: ShortTermMemoryConfig{
			MaxWindowSize: 10,
		},
	}
}

func (c *AgentConfig) Validate() error {
	if err := c.ShortTermMemoryConfig.Validate(); err != nil {
		return err
	}
	return nil
}
