package llm

import (
	"context"
	"fmt"

	"github.com/Mirai3103/Project-Re-ENE/config"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

type LLMProvider interface {
	GetModel(ctx context.Context) (*genkit.Genkit, ai.ModelArg, error)
}

type llmProvider struct {
	cfg *config.Config
}

func NewProvider(ctx context.Context, cfg *config.Config) LLMProvider {
	return &llmProvider{cfg: cfg}
}

func (p *llmProvider) GetModel(ctx context.Context) (*genkit.Genkit, ai.ModelArg, error) {
	fmt.Println("Getting model", p.cfg.LLMConfig.Provider)
	switch p.cfg.LLMConfig.Provider {
	case "gemini":
		return newGeminiModel(ctx, p.cfg.LLMConfig.GeminiConfig)
	case "openai":
		return newOpenAIModel(ctx, p.cfg.LLMConfig.OpenAIConfig)
	default:
		return nil, nil, fmt.Errorf("provider not found")
	}
}
