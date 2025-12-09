package agent

import (
	"context"
	"strings"

	"github.com/Mirai3103/Project-Re-ENE/package/utils"
	"github.com/Mirai3103/Project-Re-ENE/store"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
)

type NewFact struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}
type Memory struct {
	Content    string   `json:"content"`
	Importance float64  `json:"importance"`
	Confidence float64  `json:"confidence"`
	Tags       []string `json:"tags"`
}
type ExtractOutput struct {
	NewUserFacts      []NewFact `json:"new_user_facts"`
	NewCharacterFacts []NewFact `json:"new_character_facts"`
	Memories          []Memory  `json:"memories"`
}

type ExtractInput struct {
	ChatHistory    []*ai.Message
	UserFacts      []store.UserFact
	CharacterFacts []store.CharacterFact
	User           *store.User
	Character      *store.Character
}
type ExtractMemoryFlow = core.Flow[ExtractInput, ExtractOutput, struct{}]

func NewExtractMemoryFlow(g *genkit.Genkit, m ai.ModelArg, embeddingService *EmbeddingService) *core.Flow[ExtractInput, ExtractOutput, struct{}] {
	return genkit.DefineFlow(
		g,
		"extractMemoryFlow",
		func(ctx context.Context, in ExtractInput) (ExtractOutput, error) {
			conversationText := ConversationToText(in.ChatHistory)
			resp, err := genkit.Generate(ctx, g,
				ai.WithPrompt(conversationText),
				ai.WithOutputType(ExtractOutput{}),
				ai.WithSystem(NewExtractPrompt(in.UserFacts, in.CharacterFacts, in.User, in.Character, in.ChatHistory)),
				ai.WithModel(m),
			)
			if err != nil {
				return ExtractOutput{}, err
			}
			var extractOutput ExtractOutput
			if err := resp.Output(&extractOutput); err != nil {
				return ExtractOutput{}, err
			}
			for _, fact := range extractOutput.Memories {
				_ = embeddingService.AddMemory(ctx, &store.Memory{
					Content:    utils.Ptr(fact.Content),
					Importance: utils.Ptr(fact.Importance),
					Confidence: utils.Ptr(fact.Confidence),
					Tags:       utils.Ptr(strings.Join(fact.Tags, ",")),
				})

			}
			return extractOutput, nil
		},
	)

}

type SummaryFlow = core.Flow[ExtractInput, string, struct{}]

func NewGenSummaryFlow(g *genkit.Genkit, m ai.ModelArg) *core.Flow[ExtractInput, string, struct{}] {
	return genkit.DefineFlow(
		g,
		"genSummaryFlow",
		func(ctx context.Context, in ExtractInput) (string, error) {
			conversationText := ConversationToText(in.ChatHistory)
			resp, err := genkit.Generate(ctx, g,
				ai.WithPrompt(conversationText),
				ai.WithSystem(NewSummaryPrompt(in.Character, in.User, in.ChatHistory)),
				ai.WithModel(m),
			)
			if err != nil {
				return "", err
			}
			return resp.Text(), nil
		},
	)
}
