package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Mirai3103/Project-Re-ENE/asr"
	"github.com/Mirai3103/Project-Re-ENE/config"
	"github.com/Mirai3103/Project-Re-ENE/package/utils"
	"github.com/Mirai3103/Project-Re-ENE/store"
	"github.com/Mirai3103/Project-Re-ENE/tts"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/core/api"
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
}

type Agent struct {
	llmModel          api.Plugin
	ttsAgent          tts.TTSAgent
	asrAgent          asr.ASRAgent
	characterStore    *store.CharacterStore
	userStore         *store.UserStore
	conversationStore *store.ConversationStore
	flow              *core.Flow[FlowInput, string, string]
	g                 *genkit.Genkit
	agentConfig       *config.AgentConfig
	logger            *slog.Logger
}

func NewAgent(llmModel api.Plugin, ttsAgent tts.TTSAgent, asrAgent asr.ASRAgent, characterStore *store.CharacterStore, userStore *store.UserStore, conversationStore *store.ConversationStore, agentConfig *config.AgentConfig, logger *slog.Logger) *Agent {
	return &Agent{
		llmModel:          llmModel,
		ttsAgent:          ttsAgent,
		asrAgent:          asrAgent,
		characterStore:    characterStore,
		userStore:         userStore,
		conversationStore: conversationStore,
		agentConfig:       agentConfig,
		logger:            logger,
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
			a.g,
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
	a.g = genkit.Init(ctx, genkit.WithPlugins(a.llmModel), genkit.WithDefaultModel("googleai/gemini-2.5-flash"))
	tools, err := a.getTools(ctx)
	if err != nil {
		return err
	}
	var toolsRefs []ai.ToolRef
	for _, tool := range tools {
		toolsRefs = append(toolsRefs, tool)
	}
	agentFlow := genkit.DefineStreamingFlow(
		a.g,
		"agentFlow",
		func(ctx context.Context, input FlowInput, callback core.StreamCallback[string]) (string, error) {
			a.conversationStore.CreateConversationIfNotExists(input.ConversationID, a.agentConfig.ShortTermMemoryConfig.MaxWindowSize, input.CharacterID, input.UserID)

			messages, err := a.conversationStore.GetMessages(input.ConversationID)
			if err != nil {
				a.logger.Error("Lỗi khi lấy tin nhắn", "error", err)

			}
			historyMessages := make([]*ai.Message, len(messages))
			for i, message := range messages {
				var hm ai.Message
				err := json.Unmarshal([]byte(message.Content), &hm)
				if err != nil {
					a.logger.Error("Lỗi khi unmarshal tin nhắn", "error", err)
				} else {
					historyMessages[i] = &hm
				}
			}

			ctx = context.WithValue(ctx, ConversationID, input.ConversationID)
			ctx = context.WithValue(ctx, CharacterID, input.CharacterID)
			ctx = context.WithValue(ctx, UserID, input.UserID)
			finalResp, err := genkit.Generate(
				ctx,
				a.g,
				ai.WithSystem(NewPrompt(a.characterStore, a.userStore, a.conversationStore, &input)),
				ai.WithMessages(historyMessages...),
				ai.WithPrompt(input.Text),
				ai.WithTools(toolsRefs...),
				ai.WithStreaming(func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
					a.logger.Info("Chunk", "text", chunk.Text())
					input.chunkChan <- chunk.Text()
					return callback(ctx, chunk.Text())
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

func (a *Agent) InferSpeak(ctx context.Context, input *FlowInput) (chan SpeakResponse, error) {
	input, err := a.preProcessInput(ctx, input)
	if err != nil {
		return nil, err
	}
	input.chunkChan = make(chan string, 20)
	resultChan := make(chan SpeakResponse, 20)

	// Start flow in goroutine
	go func() {
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
