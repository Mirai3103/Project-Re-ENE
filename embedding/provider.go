package embedding

import (
	"context"
	"fmt"

	"github.com/Mirai3103/Project-Re-ENE/config"
)

type Model interface {
	Get(ctx context.Context, text string) ([]float32, error)
	Gets(ctx context.Context, texts []string) ([][]float32, error)
}

func New(ctx context.Context, cfg *config.Config) (Model, error) {
	switch cfg.EmbeddingConfig.Provider {
	case "google":
		return newGoogleGeminiModel(ctx, cfg.EmbeddingConfig.Google)
	default:
		return nil, fmt.Errorf("embedding provider not found")
	}
}
