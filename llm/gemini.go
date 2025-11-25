package llm

import (
	"context"

	"github.com/Mirai3103/Project-Re-ENE/package/utils"

	llmConfig "github.com/Mirai3103/Project-Re-ENE/config/llm"

	"github.com/cloudwego/eino-ext/components/model/gemini"
	"github.com/cloudwego/eino/components/model"
	"google.golang.org/genai"
)

func newGeminiModel(ctx context.Context, cfg *llmConfig.GeminiConfig) (model.BaseChatModel, error) {
	client, _ := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: cfg.APIKey,
	})
	model, err := gemini.NewChatModel(ctx, &gemini.Config{
		Client:      client,
		Model:       cfg.Model,
		Temperature: utils.Ptr(cfg.Temperature),
	})
	if err != nil {
		return nil, err
	}
	return model, nil
}
