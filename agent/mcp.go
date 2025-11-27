package agent

import (
	"encoding/json"
	"errors"
	"os"
)

type MCPConfig struct {
	Command *string            `json:"command"`
	Args    *[]string          `json:"args"`
	Env     *map[string]string `json:"env"`
	Enable  *bool              `json:"enable" default:"true"`
	Url     *string            `json:"url"`
}

func (c *MCPConfig) IsValid() bool {
	return c.Command != nil || c.Url != nil
}

type MCPType string

const (
	MCPTypeStdio MCPType = "stdio"
	MCPTypeSSE   MCPType = "sse"
)

func (c *MCPConfig) GetType() MCPType {
	if c.Command != nil {
		return MCPTypeStdio
	}
	return MCPTypeSSE
}

type MCPConfigFile struct {
	McpServers map[string]MCPConfig `json:"mcpServers"`
}

func ParseMCPConfigFile(mcpConfigPath string) (*MCPConfigFile, error) {
	data, err := os.ReadFile(mcpConfigPath)
	if err != nil {
		return nil, err
	}
	var mcpConfig MCPConfigFile
	err = json.Unmarshal(data, &mcpConfig)
	if err != nil {
		return nil, err
	}
	// ignore enable false and invalid config
	for name, config := range mcpConfig.McpServers {
		if !config.IsValid() {
			delete(mcpConfig.McpServers, name)
		}
		if !*config.Enable {
			delete(mcpConfig.McpServers, name)
		}
	}
	if len(mcpConfig.McpServers) == 0 {
		return nil, errors.New("no valid mcp servers found")
	}
	return &mcpConfig, nil
}
