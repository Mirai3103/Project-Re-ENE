package agent

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/plugins/mcp"
)

type MCPConfig struct {
	Command *string   `json:"command"`
	Args    *[]string `json:"args"`
	Env     *[]string `json:"env"`
	Enable  *bool     `json:"enable" default:"true"`
	Url     *string   `json:"url"`
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
func (a *Agent) parseMcpTools(ctx context.Context) ([]ai.Tool, error) {
	if !a.agentConfig.ToolsConfig.MCP.Enable {
		return []ai.Tool{}, nil
	}
	mcpConfigs, err := ParseMCPConfigFile(a.agentConfig.ToolsConfig.MCP.ConfigPath)
	if err != nil {
		return nil, err
	}
	var tools []ai.Tool
	for name, config := range mcpConfigs.McpServers {
		if config.GetType() == MCPTypeStdio {
			client, err := mcp.NewGenkitMCPClient(mcp.MCPClientOptions{
				Name: name,
				Stdio: &mcp.StdioConfig{
					Command: *config.Command,
					Args:    *config.Args,
					Env:     *config.Env,
				},
			})
			if err != nil {
				return nil, err
			}
			mcpTools, err := client.GetActiveTools(ctx, a.g)
			if err != nil {
				return nil, err
			}
			tools = append(tools, mcpTools...)
		} else if config.GetType() == MCPTypeSSE {
			client, err := mcp.NewGenkitMCPClient(mcp.MCPClientOptions{
				Name: name,
				SSE: &mcp.SSEConfig{
					BaseURL: *config.Url,
				},
			})
			if err != nil {
				return nil, err
			}
			mcpTools, err := client.GetActiveTools(ctx, a.g)
			if err != nil {
				return nil, err
			}
			tools = append(tools, mcpTools...)
		}
	}
	return tools, nil
}
