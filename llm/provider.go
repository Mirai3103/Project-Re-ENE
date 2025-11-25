package llm

import (
	"context"
	"fmt"

	"github.com/Mirai3103/Project-Re-ENE/config"

	"github.com/cloudwego/eino/components/model"
)

type LLMProvider interface {
	GetModel(ctx context.Context) (model.BaseChatModel, error)
}

type llmProvider struct {
	cfg *config.Config
}

func NewProvider(ctx context.Context, cfg *config.Config) LLMProvider {
	return &llmProvider{cfg: cfg}
}

func (p *llmProvider) GetModel(ctx context.Context) (model.BaseChatModel, error) {
	switch p.cfg.LLMConfig.Provider {
	case "gemini":
		return newGeminiModel(ctx, p.cfg.LLMConfig.GeminiConfig)
	default:
		return nil, fmt.Errorf("provider not found")
	}
}
