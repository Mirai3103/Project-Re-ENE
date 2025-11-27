package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/Mirai3103/Project-Re-ENE/agent/react"
	"github.com/Mirai3103/Project-Re-ENE/asr"
	"github.com/Mirai3103/Project-Re-ENE/config"
	"github.com/Mirai3103/Project-Re-ENE/package/utils"
	"github.com/Mirai3103/Project-Re-ENE/store"
	"github.com/Mirai3103/Project-Re-ENE/tts"

	"github.com/cloudwego/eino-ext/components/tool/googlesearch"
	einoMcp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/schema"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
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
	Stream(ctx context.Context, conversationID string, input UserInput) (chan SpeakResponse, error)
}

type assistantAgent struct {
	cfg               *config.Config
	llmModel          model.BaseChatModel
	agent             *react.Agent
	ttsAgent          tts.TTSAgent
	asrAgent          asr.ASRAgent
	logger            *slog.Logger
	characterStore    *store.CharacterStore
	userStore         *store.UserStore
	conversationStore *store.ConversationStore
	promptCache       map[string]*prompt.DefaultChatTemplate
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
	logger *slog.Logger,
	characterStore *store.CharacterStore,
	userStore *store.UserStore,
	conversationStore *store.ConversationStore,
) (AssistantAgent, error) {
	if ttsAgent == nil {
		return nil, ErrNoTTSAgent
	}
	if asrAgent == nil {
		return nil, ErrNoASRAgent
	}

	return &assistantAgent{
		cfg:               cfg,
		llmModel:          llmModel,
		ttsAgent:          ttsAgent,
		asrAgent:          asrAgent,
		logger:            logger,
		characterStore:    characterStore,
		userStore:         userStore,
		conversationStore: conversationStore,
		promptCache:       make(map[string]*prompt.DefaultChatTemplate),
	}, nil
}

func (a *assistantAgent) getSystemPrompt(characterID string, userID string) *prompt.DefaultChatTemplate {
	cacheKey := characterID + "_" + userID
	if template, ok := a.promptCache[cacheKey]; ok {
		return template
	}
	template := NewSystemPromptBuilder(a.characterStore, a.userStore).WithCharacterId(characterID).WithUserId(userID).Build()

	a.promptCache[cacheKey] = template
	return template
}

func (a *assistantAgent) mcpConfigsToTools(ctx context.Context, path string) ([]tool.BaseTool, error) {
	cfg, err := ParseMCPConfigFile(path)
	if err != nil {
		return nil, err
	}

	var tools []tool.BaseTool

	for name, server := range cfg.McpServers {
		cli, err := initMCPClient(ctx, name, server)
		if err != nil {
			continue
		}

		ts, err := einoMcp.GetTools(ctx, &einoMcp.Config{Cli: cli})
		if err != nil {
			continue
		}

		tools = append(tools, ts...)
	}

	return tools, nil
}

// -------------------------------------------------------------

func initMCPClient(ctx context.Context, name string, server MCPConfig) (*client.Client, error) {
	var (
		cli *client.Client
		err error
	)

	switch server.GetType() {

	case MCPTypeStdio:
		env := make([]string, 0)
		if server.Env != nil {
			for k, v := range *server.Env {
				env = append(env, fmt.Sprintf("%s=%s", k, v))
			}
		}

		args := []string{}
		if server.Args != nil {
			args = *server.Args
		}

		cli, err = client.NewStdioMCPClient(*server.Command, env, args...)
		if err != nil {
			return nil, err
		}

	case MCPTypeSSE:
		cli, err = client.NewSSEMCPClient(*server.Url)
		if err != nil {
			return nil, err
		}
		if err := cli.Start(ctx); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported MCP type: %s", server.GetType())
	}

	initReq := mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{
				Name:    name,
				Version: "1.0.0",
			},
		},
	}

	_, err = cli.Initialize(ctx, initReq)
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func (a *assistantAgent) CompileChain(ctx context.Context) error {
	googleTool, _ := googlesearch.NewTool(ctx, &googlesearch.Config{
		APIKey:         a.cfg.AgentConfig.ToolsConfig.GoogleSearch.APIKey,
		SearchEngineID: a.cfg.AgentConfig.ToolsConfig.GoogleSearch.SearchEngineID,
		BaseURL:        a.cfg.AgentConfig.ToolsConfig.GoogleSearch.BaseURL,
		Num:            a.cfg.AgentConfig.ToolsConfig.GoogleSearch.Num,
		Lang:           a.cfg.AgentConfig.ToolsConfig.GoogleSearch.Lang,
		ToolName:       "google_search",
		ToolDesc:       "google search tool",
	})

	mcpTools, err := a.mcpConfigsToTools(ctx, a.cfg.AgentConfig.ToolsConfig.MCP.ConfigPath)
	if err != nil {
		a.logger.Error("mcpConfigsToTools failed", "error", err)
	}
	baseTools := []tool.BaseTool{
		googleTool,
	}
	for _, tool := range mcpTools {
		info, err := tool.Info(ctx)
		if err != nil {
			continue
		}
		a.logger.Info("mcp tool", "tool", info)
	}
	tools := compose.ToolsNodeConfig{
		Tools: append(baseTools, mcpTools...),
	}
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: a.llmModel.(model.ToolCallingChatModel),
		ToolsConfig:      tools,
		MaxStep:          20,
		StreamToolCallChecker: func(ctx context.Context, modelOutput *schema.StreamReader[*schema.Message]) (bool, error) {
			var isToolCall = false
			cpyModelOutput := modelOutput.Copy(1)[0]
			for {
				part, err := cpyModelOutput.Recv()
				if err == io.EOF || errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					return false, err
				}
				if len(part.ToolCalls) > 0 {
					isToolCall = true
				}
			}
			return isToolCall, nil
		},
	})
	if err != nil {
		return err
	}
	a.agent = agent
	return nil
}
func (a *assistantAgent) createConversationIfNotExist(conversationID string, characterID string, userID string) error {
	_, err := a.conversationStore.GetConversation(conversationID)
	if errors.Is(err, store.ErrNoRecordFound) {
		err = a.conversationStore.CreateConversation(conversationID, 10, characterID, userID)
		if err != nil {
			a.logger.Error("createConversationIfNotExist failed", "error", err)
			return err
		}
	}
	return nil
}
func (a *assistantAgent) getChatHistory(conversationID string) ([]*schema.Message, error) {
	messages, err := a.conversationStore.GetMessages(conversationID)
	if err != nil {
		return nil, err
	}
	chatHistory := make([]*schema.Message, len(messages))
	for i, message := range messages {
		jsonRaw := []byte(message.Content)
		var message schema.Message
		err = json.Unmarshal(jsonRaw, &message)
		if err != nil {
			return nil, err
		}
		chatHistory[i] = &message

	}
	return chatHistory, nil
}
func (a *assistantAgent) appendMessage(conversationID string, message *schema.Message) error {
	role := string(message.Role)
	jsonRaw, err := json.Marshal(message)
	if err != nil {
		return err
	}
	content := string(jsonRaw)
	err = a.conversationStore.AppendMessage(conversationID, &store.ConversationMessage{
		Role:      role,
		Content:   content,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *assistantAgent) Stream(ctx context.Context, conversationID string, input UserInput) (chan SpeakResponse, error) {
	log := a.logger
	log.Info("Streaming input", "input", input)
	inputText, err := a.processInput(ctx, input)
	if err != nil {
		return nil, err
	}

	chatHistory, err := a.getChatHistory(conversationID)
	if err != nil {
		return nil, err
	}
	prompts, err := a.getSystemPrompt("1", "huuhoang").Format(ctx, map[string]any{
		USER_INPUT_KEY: inputText,
		USER_IMAGE_KEY: input.Image,
		"chat_history": chatHistory,
	})
	if err != nil {
		return nil, err
	}
	contentChan := make(chan string, 100)
	cb := callbacks.NewHandlerBuilder().OnEndWithStreamOutputFn(
		func(ctx context.Context, info *callbacks.RunInfo, output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {
			cloneOutput := output.Copy(1)[0]
			for {
				part, err := cloneOutput.Recv()
				if err == io.EOF || errors.Is(err, io.EOF) {
					break
				}

				if v, ok := part.(*model.CallbackOutput); ok {
					contentChan <- v.Message.Content
				}
			}
			return ctx
		}).Build()
	ch := make(chan SpeakResponse)
	go a.handleStream(ctx, conversationID, inputText, ch, contentChan)
	go func() {
		defer close(contentChan)
		stream, err := a.agent.Stream(ctx, prompts, agent.WithComposeOptions(compose.WithCallbacks(cb)))
		if err != nil {
			return
		}
		// defer close(stringOutPutChan) // Đóng channel khi stream kết thúc
		for {
			_, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					log.Info("Stream source finished (EOF)")
				} else {
					log.Error("Stream source error", "error", err)
				}
				return
			}
		}
	}()

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
	ch chan SpeakResponse,
	stringOutPutChan chan string,

) {
	log := a.logger
	defer close(ch)

	var fullResponse strings.Builder
	buffer := ""

	// 1. Goroutine phụ: Dùng để duy trì stream và đóng channel text khi xong

	// Always save conversation before returning
	defer func() {
		err := a.appendMessage(conversationID, schema.AssistantMessage(fullResponse.String(), nil))
		if err != nil {
			log.Error("appendMessage failed", "error", err)
			return
		}
		err = a.appendMessage(conversationID, schema.UserMessage(inputText))
		if err != nil {
			log.Error("appendMessage failed", "error", err)
			return
		}
	}()

	// 2. Luồng chính: Xử lý text từ channel
	for content := range stringOutPutChan {
		// Check context cancellation
		if err := ctx.Err(); err != nil {
			log.Error("Context cancelled", "error", err, "fullResponse", fullResponse.String())
			a.flushBuffer(ctx, buffer, ch, log)
			return
		}

		// Accumulate full response
		fullResponse.WriteString(content)
		buffer += content

		// Process complete sentences
		buffer = a.processAndSendSentences(ctx, buffer, ch, log)
	}

	// 3. Khi stringOutPutChan đóng (vòng for ở trên kết thúc), flush phần còn lại
	log.Info("Channel closed, flushing buffer", "fullResponse", fullResponse.String())
	a.flushBuffer(ctx, buffer, ch, log)
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
