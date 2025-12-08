package llm

import (
	"context"

	llmConfig "github.com/Mirai3103/Project-Re-ENE/config/llm"
	"github.com/openai/openai-go/option"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	oai "github.com/firebase/genkit/go/plugins/compat_oai"
)

const key = "random"

func newOpenAIModel(ctx context.Context, cfg *llmConfig.OpenAIConfig) (*genkit.Genkit, ai.ModelArg, error) {
	o := &oai.OpenAICompatible{
		APIKey: cfg.APIKey,
		Opts: []option.RequestOption{
			option.WithAPIKey(cfg.APIKey),
			option.WithBaseURL(cfg.BaseURL),
		},
		Provider: key,
		BaseURL:  cfg.BaseURL,
	}
	g := genkit.Init(ctx, genkit.WithPlugins(o), genkit.WithDefaultModel(key+"/"+cfg.Model))
	modelRef := o.DefineModel("", key+"/"+cfg.Model, ai.ModelOptions{
		Supports: &ai.ModelSupports{
			Tools:      true,
			SystemRole: true,
		},
	})
	return g, modelRef, nil
}
