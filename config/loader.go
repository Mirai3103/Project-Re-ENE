package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func GetDefaultConfig() *Config {
	return &Config{
		LLMConfig:       *getDefaultLLMConfig(),
		TTSConfig:       *getDefaultTTSConfig(),
		LoggerConfig:    *getDefaultLoggerConfig(),
		CharacterConfig: *getDefaultCharacterConfig(),
		AgentConfig:     *getDefaultAgentConfig(),
		ModelsConfig:    *getDefaultModelsConfig(),
	}
}

func LoadConfig(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := GetDefaultConfig()
		data, err := yaml.Marshal(defaultConfig)
		if err != nil {
			return nil, fmt.Errorf("marshal default config: %w", err)
		}
		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return nil, fmt.Errorf("write default config: %w", err)
		}
		fmt.Println("Created default config file:", configPath)
		return defaultConfig, nil
	}

	// nếu file có → đọc
	fmt.Println("Loading config from:", configPath)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, cfg.Validate()
}

func PersistConfig(cfg *Config, configPath string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}
