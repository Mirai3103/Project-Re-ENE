package llm

import (
	"context"
	"fmt"

	"github.com/Mirai3103/Project-Re-ENE/config"
	"github.com/firebase/genkit/go/core/api"
)

type LLMProvider interface {
	GetModel(ctx context.Context) (api.Plugin, error)
}

type llmProvider struct {
	cfg *config.Config
}

func NewProvider(ctx context.Context, cfg *config.Config) LLMProvider {
	return &llmProvider{cfg: cfg}
}

func (p *llmProvider) GetModel(ctx context.Context) (api.Plugin, error) {
	switch p.cfg.LLMConfig.Provider {
	case "gemini":
		return newGeminiModel(ctx, p.cfg.LLMConfig.GeminiConfig), nil
	default:
		return nil, fmt.Errorf("provider not found")
	}
}
