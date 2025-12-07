package embedding

import (
	"context"

	"github.com/Mirai3103/Project-Re-ENE/config/embedding"
	"google.golang.org/genai"
)

type googleGeminiModel struct {
	cfg    *embedding.GoogleGeminiConfig
	client *genai.Client
}

func NewGoogleGeminiModel(ctx context.Context, cfg *embedding.GoogleGeminiConfig) (Model, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: cfg.APIKey,
	})
	if err != nil {
		return nil, err
	}
	return &googleGeminiModel{cfg: cfg, client: client}, nil
}

func (m *googleGeminiModel) Get(ctx context.Context, text string) ([]float32, error) {
	contents := []*genai.Content{
		genai.NewContentFromText(text, genai.RoleUser),
	}
	result, err := m.client.Models.EmbedContent(ctx,
		m.cfg.ModelID,
		contents,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return result.Embeddings[0].Values, nil
}

func (m *googleGeminiModel) Gets(ctx context.Context, texts []string) ([][]float32, error) {
	contents := make([]*genai.Content, len(texts))
	for i, text := range texts {
		contents[i] = genai.NewContentFromText(text, genai.RoleUser)
	}
	result, err := m.client.Models.EmbedContent(ctx,
		m.cfg.ModelID,
		contents,
		nil,
	)
	if err != nil {
		return nil, err
	}
	embeddings := make([][]float32, len(result.Embeddings))
	for i, embedding := range result.Embeddings {
		embeddings[i] = embedding.Values
	}
	return embeddings, nil
}
