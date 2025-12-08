package llm

import (
	"context"
	"fmt"

	llmConfig "github.com/Mirai3103/Project-Re-ENE/config/llm"
	"google.golang.org/genai"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
)

func newGeminiModel(ctx context.Context, cfg *llmConfig.GeminiConfig) (*genkit.Genkit, ai.ModelArg, error) {
	gemini := &googlegenai.GoogleAI{
		APIKey: cfg.APIKey,
	}
	g := genkit.Init(ctx, genkit.WithPlugins(gemini), genkit.WithDefaultModel(cfg.Model))
	modelRef := googlegenai.GoogleAIModelRef("gemini-2.5-flash", &genai.GenerateContentConfig{
		Temperature: genai.Ptr[float32](cfg.Temperature),
		// Other configuration...
	})
	fmt.Println("Gemini model initialized", cfg.Model)
	return g, modelRef, nil
}
