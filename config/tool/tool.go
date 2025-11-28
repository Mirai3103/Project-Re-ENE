package tools

type ToolConfig struct {
	GoogleSearch GoogleSearchToolConfig `yaml:"google_search"`
	MCP          MCPToolConfig          `yaml:"mcp"`
}

type GoogleSearchToolConfig struct {
	APIKey         string `yaml:"api_key"`
	SearchEngineID string `yaml:"search_engine_id"`
	BaseURL        string `yaml:"base_url"`
	Num            int    `yaml:"num"`
	Lang           string `yaml:"lang"`
	Enable         bool   `yaml:"enable"`
}

type MCPToolConfig struct {
	ConfigPath string `yaml:"config_path"`
	Enable     bool   `yaml:"enable"`
}

func getDefaultMCPToolConfig() *MCPToolConfig {
	return &MCPToolConfig{
		ConfigPath: "./resources/mcp/config.json",
		Enable:     false,
	}
}
func getDefaultGoogleSearchToolConfig() *GoogleSearchToolConfig {
	return &GoogleSearchToolConfig{
		APIKey:         "",
		SearchEngineID: "",
		BaseURL:        "https://customsearch.googleapis.com",
		Num:            5,
		Lang:           "vi",
		Enable:         false,
	}
}

func GetDefaultToolConfig() *ToolConfig {
	return &ToolConfig{
		GoogleSearch: *getDefaultGoogleSearchToolConfig(),
		MCP:          *getDefaultMCPToolConfig(),
	}
}
