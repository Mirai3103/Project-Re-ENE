package config

import (
	"errors"
	"fmt"
	"os"
)

type ShortTermMemoryConfig struct {
	MaxWindowSize    int    `yaml:"max_window_size"`
	ConversationsDir string `yaml:"conversation_dir"`
}

func (c *ShortTermMemoryConfig) Validate() error {
	if c.MaxWindowSize <= 0 {
		return errors.New("max_window_size must be greater than 0")
	}
	if c.ConversationsDir == "" {
		return errors.New("conversation_dir must be set")
	}
	if err := os.MkdirAll(c.ConversationsDir, 0755); err != nil {
		return errors.New("failed to create conversations dir")
	}
	fmt.Println("ConversationsDir: ", c.ConversationsDir)
	return nil
}

type AgentConfig struct {
	ShortTermMemoryConfig ShortTermMemoryConfig `yaml:"short_term_memory_config"`
}

func getDefaultAgentConfig() *AgentConfig {
	return &AgentConfig{
		ShortTermMemoryConfig: ShortTermMemoryConfig{
			MaxWindowSize:    10,
			ConversationsDir: ".data/conversations",
		},
	}
}

func (c *AgentConfig) Validate() error {
	if err := c.ShortTermMemoryConfig.Validate(); err != nil {
		return err
	}
	return nil
}
