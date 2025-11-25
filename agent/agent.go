package agent

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"

	"github.com/Mirai3103/Project-Re-ENE/asr"
	"github.com/Mirai3103/Project-Re-ENE/config"
	"github.com/Mirai3103/Project-Re-ENE/package/utils"
	"github.com/Mirai3103/Project-Re-ENE/tts"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

var (
	ErrNoTTSAgent    = errors.New("you must provide a tts agent")
	ErrNoASRAgent    = errors.New("you must provide an asr agent")
	ErrNoInputData   = errors.New("you must provide a text or audio")
	ErrStreamRecv    = errors.New("error receiving stream data")
	ErrContextCancel = errors.New("context cancelled")
)

type AssistantAgent interface {
	CompileChain(ctx context.Context) error
	Invoke(ctx context.Context, conversationID string, input UserInput) (SpeakResponse, error)
	Stream(ctx context.Context, conversationID string, input UserInput) (chan SpeakResponse, error)
}

type assistantAgent struct {
	cfg             *config.Config
	llmModel        model.BaseChatModel
	chain           compose.Runnable[map[string]any, *schema.Message]
	ttsAgent        tts.TTSAgent
	asrAgent        asr.ASRAgent
	shortTermMemory *SimpleMemory
	logger          *slog.Logger
}

type UserInput struct {
	Text  *string `json:"text"`
	Audio *[]byte `json:"audio"`
	Image *[]byte `json:"image"`
}

func NewAssistantAgent(
	cfg *config.Config,
	llmModel model.BaseChatModel,
	ttsAgent tts.TTSAgent,
	asrAgent asr.ASRAgent,
) (AssistantAgent, error) {
	if ttsAgent == nil {
		return nil, ErrNoTTSAgent
	}
	if asrAgent == nil {
		return nil, ErrNoASRAgent
	}

	return &assistantAgent{
		cfg:             cfg,
		llmModel:        llmModel,
		ttsAgent:        ttsAgent,
		asrAgent:        asrAgent,
		shortTermMemory: NewSimpleMemory(cfg.AgentConfig.ShortTermMemoryConfig),
	}, nil
}

func (a *assistantAgent) CompileChain(ctx context.Context) error {
	chain, err := compose.NewChain[map[string]any, *schema.Message]().
		AppendChatTemplate(systemPrompt).
		AppendChatModel(a.llmModel).
		Compile(ctx)
	if err != nil {
		return err
	}
	a.chain = chain
	return nil
}

func (a *assistantAgent) Invoke(ctx context.Context, conversationID string, input UserInput) (SpeakResponse, error) {
	log := a.logger
	conversation := a.shortTermMemory.GetConversation(conversationID, true)

	inputText, err := a.processInput(ctx, input)
	log.Info("inputText: ", "inputText", inputText)
	if err != nil {
		return SpeakResponse{}, err
	}

	message, err := a.chain.Invoke(ctx, map[string]any{
		USER_INPUT_KEY: inputText,
		USER_IMAGE_KEY: input.Image,
		"chat_history": conversation.GetMessages(),
	})
	if err != nil {
		return SpeakResponse{}, err
	}

	responseText := message.Content

	go func() {
		conversation.Append(schema.UserMessage(inputText))
		conversation.Append(schema.AssistantMessage(responseText, nil))
	}()

	audioBuffer, err := a.ttsAgent.GetTTS(ctx, responseText)
	if err != nil {
		return SpeakResponse{}, err
	}

	return SpeakResponse{
		Text:        responseText,
		AudioBuffer: audioBuffer,
	}, nil
}

func (a *assistantAgent) Stream(ctx context.Context, conversationID string, input UserInput) (chan SpeakResponse, error) {
	inputText, err := a.processInput(ctx, input)
	if err != nil {
		return nil, err
	}

	conversation := a.shortTermMemory.GetConversation(conversationID, true)
	stream, err := a.chain.Stream(ctx, map[string]any{
		USER_INPUT_KEY: inputText,
		USER_IMAGE_KEY: input.Image,
		"chat_history": conversation.GetMessages(),
	})
	if err != nil {
		return nil, err
	}

	ch := make(chan SpeakResponse)

	go a.handleStream(ctx, conversationID, inputText, stream, ch)

	return ch, nil
}

// processInput extracts text from UserInput (either from Text or Audio via ASR)
func (a *assistantAgent) processInput(ctx context.Context, input UserInput) (string, error) {

	if input.Text != nil {
		return *input.Text, nil
	}

	if input.Audio != nil {
		text, err := a.asrAgent.GetASR(ctx, *input.Audio)

		if err != nil {
			return "", err
		}
		if strings.TrimSpace(text) == "" {
			return "", ErrNoInputData
		}
		return text, nil
	}

	return "", ErrNoInputData
}

// handleStream processes the LLM stream, generates TTS, and sends responses
func (a *assistantAgent) handleStream(
	ctx context.Context,
	conversationID string,
	inputText string,
	stream *schema.StreamReader[*schema.Message],
	ch chan SpeakResponse,
) {
	log := a.logger
	defer close(ch)

	var fullResponse strings.Builder
	buffer := ""

	// Always save conversation before returning
	defer func() {
		a.saveConversation(conversationID, inputText, fullResponse.String())
	}()

	for {
		// Check context cancellation
		if err := ctx.Err(); err != nil {
			log.Error("Context cancelled", "error", err, "fullResponse", fullResponse.String())
			a.flushBuffer(ctx, buffer, ch, log)
			return
		}

		part, err := stream.Recv()
		if err != nil {
			if err == io.EOF || errors.Is(err, io.EOF) {
				log.Info("Stream completed", "fullResponse", fullResponse.String())
			} else {
				log.Error("Stream error", "error", err, "fullResponse", fullResponse.String())
			}

			// Flush remaining buffer
			a.flushBuffer(ctx, buffer, ch, log)
			return
		}

		// Accumulate full response
		fullResponse.WriteString(part.Content)
		buffer += part.Content

		// Process complete sentences
		buffer = a.processAndSendSentences(ctx, buffer, ch, log)
	}
}

// processAndSendSentences splits buffer into sentences, processes complete ones, returns incomplete
func (a *assistantAgent) processAndSendSentences(
	ctx context.Context,
	buffer string,
	ch chan SpeakResponse,
	log *slog.Logger,
) string {
	sentences := utils.SplitSentences(buffer)

	// Process all complete sentences (all except the last one which might be incomplete)
	for i := 0; i < len(sentences)-1; i++ {
		sentence := strings.TrimSpace(sentences[i])
		if sentence == "" {
			continue
		}

		// Check context before TTS
		if err := ctx.Err(); err != nil {
			log.Error("Context cancelled during sentence processing", "error", err)
			return ""
		}

		audio, err := a.ttsAgent.GetTTS(ctx, sentence)
		if err != nil {
			log.Error("TTS generation failed", "error", err, "sentence", sentence)
			continue
		}

		// Send response with context check
		select {
		case <-ctx.Done():
			log.Error("Context cancelled during send", "error", ctx.Err())
			return ""
		case ch <- SpeakResponse{Text: sentence, AudioBuffer: audio}:
		}
	}

	// Return the last (potentially incomplete) sentence as the new buffer
	if len(sentences) > 0 {
		return sentences[len(sentences)-1]
	}
	return ""
}

// flushBuffer processes and sends any remaining text in the buffer
func (a *assistantAgent) flushBuffer(
	ctx context.Context,
	buffer string,
	ch chan SpeakResponse,
	log *slog.Logger,
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

		audio, err := a.ttsAgent.GetTTS(ctx, sentence)
		if err != nil {
			log.Error("TTS generation failed in flush", "error", err, "sentence", sentence)
			continue
		}

		select {
		case <-ctx.Done():
			log.Error("Context cancelled during flush", "error", ctx.Err())
			return
		case ch <- SpeakResponse{Text: sentence, AudioBuffer: audio}:
		}
	}
}

// saveConversation persists the conversation to memory
func (a *assistantAgent) saveConversation(conversationID, input, output string) {
	conversation := a.shortTermMemory.GetConversation(conversationID, true)
	conversation.Append(schema.UserMessage(input))
	conversation.Append(schema.AssistantMessage(output, nil))
}
