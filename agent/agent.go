package agent

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/Mirai3103/Project-Re-ENE/asr"
	"github.com/Mirai3103/Project-Re-ENE/config"
	"github.com/Mirai3103/Project-Re-ENE/package/utils"
	"github.com/Mirai3103/Project-Re-ENE/store"
	"github.com/Mirai3103/Project-Re-ENE/tts"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

type FlowInput struct {
	Text           string
	chunkChan      chan string
	Audio          []byte
	ConversationID string
	UserID         string
	CharacterID    string
	UserFacts      []store.UserFact
	CharacterFacts []store.CharacterFact
	User           *store.User
	Character      *store.Character
}

type Agent struct {
	llmModel          *genkit.Genkit
	ttsAgent          tts.TTSAgent
	asrAgent          asr.ASRAgent
	store             *store.Queries
	flow              *core.Flow[FlowInput, string, string]
	agentConfig       *config.AgentConfig
	logger            *slog.Logger
	modelArg          ai.ModelArg
	extractMemoryFlow *ExtractMemoryFlow
	summaryFlow       *SummaryFlow
	embeddingService  *EmbeddingService
}

func NewAgent(llmModel *genkit.Genkit, modelArg ai.ModelArg, embeddingService *EmbeddingService, ttsAgent tts.TTSAgent, asrAgent asr.ASRAgent, store *store.Queries, agentConfig *config.AgentConfig, logger *slog.Logger) *Agent {
	return &Agent{
		llmModel:          llmModel,
		ttsAgent:          ttsAgent,
		asrAgent:          asrAgent,
		store:             store,
		agentConfig:       agentConfig,
		logger:            logger,
		modelArg:          modelArg,
		extractMemoryFlow: NewExtractMemoryFlow(llmModel, modelArg, embeddingService),
		summaryFlow:       NewGenSummaryFlow(llmModel, modelArg),
		embeddingService:  embeddingService,
	}
}

type GoogleSearchInput struct {
	Query string
	Num   int64
}

func (a *Agent) getTools(ctx context.Context) ([]ai.Tool, error) {
	var tools []ai.Tool
	if a.agentConfig.ToolsConfig.GoogleSearch.Enable {
		searchService, err := customsearch.NewService(ctx, option.WithAPIKey(a.agentConfig.ToolsConfig.GoogleSearch.APIKey))
		if err != nil {
			goto afterGoogleSearch
		}
		googleSearchTool := genkit.DefineTool(
			a.llmModel,
			"googleSearch",
			"Searches the web for a given query",
			func(ctx *ai.ToolContext, input GoogleSearchInput) (any, error) {
				results, err := searchService.Cse.List().Q(input.Query).Num(input.Num).Cx(a.agentConfig.ToolsConfig.GoogleSearch.SearchEngineID).Do()
				if err != nil {
					return nil, err
				}
				return results.Items, nil
			},
		)
		tools = append(tools, googleSearchTool)
	}
afterGoogleSearch:
	mcpTools, err := a.parseMcpTools(ctx)
	if err == nil {
		tools = append(tools, mcpTools...)
	}
	for _, tool := range tools {
		a.logger.Info("Registered tool", "name", tool.Name(), "description", tool.Definition().Description)
	}
	return tools, nil
}

type ContextKey string

const (
	ConversationID ContextKey = "conversationID"
	CharacterID    ContextKey = "characterID"
	UserID         ContextKey = "userID"
)

func (a *Agent) Compile(ctx context.Context) error {
	tools, err := a.getTools(ctx)
	if err != nil {
		return err
	}
	var toolsRefs []ai.ToolRef
	for _, tool := range tools {
		toolsRefs = append(toolsRefs, tool)
	}

	agentFlow := genkit.DefineStreamingFlow(
		a.llmModel,
		"agentFlow",
		func(ctx context.Context, input FlowInput, callback core.StreamCallback[string]) (string, error) {
			cvs, err := a.store.CreateConversationIfNotExists(ctx, store.CreateConversationParams{
				ID:          input.ConversationID,
				UserID:      utils.Ptr(input.UserID),
				CharacterID: utils.Ptr(input.CharacterID),
			})
			if err != nil {
				return "", err
			}
			var windowSize = int64(a.agentConfig.ShortTermMemoryConfig.MaxWindowSize)
			if cvs.MaxWindowSize != nil {
				windowSize = *cvs.MaxWindowSize
			}
			messages, err := a.store.ListRecentMessages(ctx, store.ListRecentMessagesParams{
				ConversationID: utils.Ptr(input.ConversationID),
				Limit:          windowSize,
			})
			if err != nil {
				a.logger.Error("Lỗi khi lấy tin nhắn", "error", err)

			}
			historyMessages := ParseHistoryMessages(messages)

			ctx = context.WithValue(ctx, ConversationID, input.ConversationID)
			ctx = context.WithValue(ctx, CharacterID, input.CharacterID)
			ctx = context.WithValue(ctx, UserID, input.UserID)
			finalResp, err := genkit.Generate(
				ctx,
				a.llmModel,
				ai.WithModel(a.modelArg),
				ai.WithSystem(NewPrompt(input.UserFacts, input.CharacterFacts, input.User, input.Character)),
				ai.WithMessages(historyMessages...),
				ai.WithPrompt(input.Text),
				ai.WithTools(toolsRefs...),
				ai.WithStreaming(func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
					a.logger.Info("Chunk", "text", chunk.Text())
					trimmedText := strings.TrimSpace(chunk.Text())
					if trimmedText != "" {
						input.chunkChan <- chunk.Text()
						return callback(ctx, chunk.Text())

					}
					return nil

				}),
				ai.WithMiddleware(a.SaveConversationMiddleware),
			)
			if err != nil {
				a.logger.Error("Generation error", "error", err)
				return "", err
			}
			return finalResp.Text(), nil
		},
	)
	a.flow = agentFlow
	return nil
}

func (a *Agent) preProcessInput(ctx context.Context, input *FlowInput) (*FlowInput, error) {
	trimmedText := strings.TrimSpace(input.Text)

	// Case 1: Empty input
	if trimmedText == "" && input.Audio == nil {
		return nil, errors.New("text or audio is required")
	}

	// Case 2: Has text - prioritize text over audio
	if trimmedText != "" {
		input.Text = trimmedText
		input.Audio = nil // Clear audio to avoid confusion
		return input, nil
	}

	// Case 3: Only has audio - transcribe it
	transcribedText, err := a.asrAgent.GetASR(ctx, input.Audio)
	if err != nil {
		return nil, fmt.Errorf("ASR transcription failed: %w", err)
	}

	input.Text = strings.TrimSpace(transcribedText)

	// Case 4: Transcription returned empty text
	if input.Text == "" {
		return nil, errors.New("transcription returned empty text")
	}

	return input, nil
}
func (a *Agent) RetrieveRelatedInfo(ctx context.Context, input *FlowInput) *FlowInput {
	user, _ := a.store.GetUser(ctx, input.UserID)
	userFacts, _ := a.store.GetUserFacts(ctx, store.GetUserFactsParams{
		Limit:  10,
		UserID: utils.Ptr(input.UserID),
	})
	characterFacts, _ := a.store.GetCharacterFacts(ctx, store.GetCharacterFactsParams{
		Limit:       10,
		CharacterID: utils.Ptr(input.CharacterID),
	})

	character, _ := a.store.GetCharacter(ctx, input.CharacterID)
	input.User = &user
	input.UserFacts = userFacts
	input.CharacterFacts = characterFacts
	input.Character = &character
	return input
}

func (a *Agent) InferSpeak(ctx context.Context, input *FlowInput) (chan SpeakResponse, error) {
	input, err := a.preProcessInput(ctx, input)
	if err != nil {
		return nil, err
	}
	input = a.RetrieveRelatedInfo(ctx, input)

	input.chunkChan = make(chan string, 20)
	resultChan := make(chan SpeakResponse, 20)

	// Start flow in goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		input.Audio = []byte(input.Text)
		chunks := a.flow.Stream(ctx, *input)

		a.logger.Info("Flow completed", "final response", chunks)
		for chunk := range chunks {
			if chunk == nil {
				continue
			}
			if chunk.Done {
				break
			}
		}

		close(input.chunkChan) // Signal that streaming is complete
	}()

	// Process chunks and convert to speech
	go a.handleStreamToSpeech(ctx, input.chunkChan, resultChan)

	go func() {
		wg.Wait()
		bgCtx := context.Background()
		historyMessages, _ := a.store.ListConversationMessages(ctx, utils.Ptr(input.ConversationID))
		if len(historyMessages)%20 == 0 && len(historyMessages) > 0 {
			summary, err := a.summaryFlow.Run(bgCtx, ExtractInput{
				ChatHistory:    ParseHistoryMessages(historyMessages),
				UserFacts:      input.UserFacts,
				CharacterFacts: input.CharacterFacts,
				User:           input.User,
				Character:      input.Character})
			if err != nil {
				a.logger.Error("Lỗi khi tạo summary", "error", err)
				return
			}
			a.logger.Info("Summary", "summary", summary)
			// todo: update summary to database
		}
		if len(historyMessages) > 0 {
			facts, err := a.extractMemoryFlow.Run(bgCtx, ExtractInput{
				ChatHistory:    ParseHistoryMessages(historyMessages),
				UserFacts:      input.UserFacts,
				CharacterFacts: input.CharacterFacts,
				User:           input.User,
				Character:      input.Character})
			if err != nil {
				a.logger.Error("Lỗi khi tạo facts", "error", err)
				return
			}
			a.logger.Info("Facts", "facts", facts)
		}
		// todo: save facts to database
	}()
	return resultChan, nil
}

func (a *Agent) handleStreamToSpeech(
	ctx context.Context,
	chunkChan <-chan string,
	resultChan chan<- SpeakResponse,
) {
	defer close(resultChan)

	var buffer string

	for chunk := range chunkChan {
		// Check context cancellation
		if err := ctx.Err(); err != nil {
			a.logger.Error("Context cancelled", "error", err)
			return
		}

		buffer += chunk
		buffer = a.processCompleteSentences(ctx, buffer, resultChan)
	}

	// Process remaining buffer
	a.flushBuffer(ctx, buffer, resultChan)
}

func (a *Agent) processCompleteSentences(
	ctx context.Context,
	buffer string,
	resultChan chan<- SpeakResponse,
) string {
	sentences := utils.SplitSentences(buffer)

	// Process all complete sentences (all except last which might be incomplete)
	for i := 0; i < len(sentences)-1; i++ {
		sentence := strings.TrimSpace(sentences[i])
		if sentence == "" {
			continue
		}

		if err := a.sendSpeechResponse(ctx, sentence, resultChan); err != nil {
			return "" // Stop processing on context cancellation
		}
	}

	// Return the last (potentially incomplete) sentence
	if len(sentences) > 0 {
		return sentences[len(sentences)-1]
	}
	return ""
}

func (a *Agent) flushBuffer(
	ctx context.Context,
	buffer string,
	resultChan chan<- SpeakResponse,
) {
	if buffer == "" {
		return
	}

	sentences := utils.SplitSentences(buffer)
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence == "" {
			continue
		}

		if err := a.sendSpeechResponse(ctx, sentence, resultChan); err != nil {
			return // Stop on context cancellation
		}
	}
}

func (a *Agent) sendSpeechResponse(
	ctx context.Context,
	text string,
	resultChan chan<- SpeakResponse,
) error {
	// Check context before expensive TTS operation
	if err := ctx.Err(); err != nil {
		a.logger.Error("Context cancelled before TTS", "error", err)
		return err
	}

	audio, err := a.ttsAgent.GetTTS(ctx, text)
	if err != nil {
		a.logger.Error("TTS generation failed", "error", err, "text", text)
		return nil // Continue processing other sentences
	}

	// Send with context check
	select {
	case <-ctx.Done():
		a.logger.Error("Context cancelled during send", "error", ctx.Err())
		return ctx.Err()
	case resultChan <- SpeakResponse{Text: text, AudioBuffer: audio}:
		return nil
	}
}
