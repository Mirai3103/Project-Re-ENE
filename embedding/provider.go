package embedding

import (
	"context"

	"github.com/Mirai3103/Project-Re-ENE/config"
)

type Model interface {
	Get(ctx context.Context, text string) ([]float32, error)
	Gets(ctx context.Context, texts []string) ([][]float32, error)
}

type Provider interface {
	GetModel(ctx context.Context) (Model, error)
}

type provider struct {
	cfg *config.Config
}
