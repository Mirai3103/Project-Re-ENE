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
	Name  string `json:"name"  jsonschema:"description=Short title representing the fact being extracted"`
	Value string `json:"value" jsonschema:"description=The actual factual information stated clearly in the conversation"`
	Type  string `json:"type"  jsonschema:"description=Fact category: 'user' or 'character'"`
}

type Memory struct {
	Content    string   `json:"content"    jsonschema:"description=Diary-like narrative describing the event or emotion from the Character's POV"`
	Importance float64  `json:"importance" jsonschema:"description=Importance score from 0.0 to 1.0 based on emotional weight or long-term relevance"`
	Confidence float64  `json:"confidence" jsonschema:"description=How certain the system is that this memory is accurate and grounded in the conversation"`
	Tags       []string `json:"tags"       jsonschema:"description=Short keywords summarizing emotion, topic, or action for easy retrieval"`
}

type ExtractOutput struct {
	// NewUserFacts      []NewFact `json:"new_user_facts"      jsonschema:"description=New factual information about the User extracted from the conversation, must be empty if no new information is found"`
	// NewCharacterFacts []NewFact `json:"new_character_facts" jsonschema:"description=New stable facts about the Character derived from their behavior or dialogue patterns, must be empty if no new information is found"`
	Memories []Memory `json:"memories"            jsonschema:"description=Significant events or emotional details worth storing for future interactions must be empty if no new information is found"`
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
			// for _, fact := range extractOutput.NewUserFacts {
			// 	embeddingService.AddUserFact(ctx, &store.UserFact{
			// 		UserID: utils.Ptr(in.User.ID),
			// 		Name:   utils.Ptr(fact.Name),
			// 		Value:  utils.Ptr(fact.Value),
			// 		Type:   utils.Ptr(fact.Type),
			// 	})
			// }
			// for _, fact := range extractOutput.NewCharacterFacts {
			// 	embeddingService.AddCharacterFact(ctx, &store.CharacterFact{
			// 		CharacterID: utils.Ptr(in.Character.ID),
			// 		Name:        utils.Ptr(fact.Name),
			// 		Value:       utils.Ptr(fact.Value),
			// 		Type:        utils.Ptr(fact.Type),
			// 	})
			// }

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
