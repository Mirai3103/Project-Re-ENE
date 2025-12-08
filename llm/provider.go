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

func New(ctx context.Context, cfg *config.Config) (*genkit.Genkit, ai.ModelArg, error) {
	fmt.Println("Getting model", cfg.LLMConfig.Provider)
	switch cfg.LLMConfig.Provider {
	case "gemini":
		return newGeminiModel(ctx, cfg.LLMConfig.GeminiConfig)
	case "openai":
		return newOpenAIModel(ctx, cfg.LLMConfig.OpenAIConfig)
	default:
		return nil, nil, fmt.Errorf("provider not found")
	}
}
