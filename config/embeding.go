package config

import "github.com/Mirai3103/Project-Re-ENE/config/embedding"

type EmbeddingConfig struct {
	Provider string                           `yaml:"provider"`
	Google   *embedding.GoogleEmbeddingConfig `yaml:"google"`
}
