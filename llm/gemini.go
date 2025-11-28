package llm

import (
	"context"

	llmConfig "github.com/Mirai3103/Project-Re-ENE/config/llm"

	"github.com/firebase/genkit/go/core/api"
	"github.com/firebase/genkit/go/plugins/googlegenai"
)

func newGeminiModel(ctx context.Context, cfg *llmConfig.GeminiConfig) api.Plugin {
	return &googlegenai.GoogleAI{
		APIKey: cfg.APIKey,
	}
}
