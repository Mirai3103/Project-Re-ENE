package agent_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/Mirai3103/Project-Re-ENE/agent"
	"github.com/Mirai3103/Project-Re-ENE/asr"
	"github.com/Mirai3103/Project-Re-ENE/config"
	"github.com/Mirai3103/Project-Re-ENE/tts"
	"github.com/firebase/genkit/go/plugins/googlegenai"
)

type MockTTS struct{}

func (m MockTTS) GetTTS(ctx context.Context, text string) ([]byte, error) {
	return []byte("mock audio data"), nil
}

func NewMockTTS() tts.TTSAgent {
	return &MockTTS{}
}

type MockASR struct{}

func (m MockASR) GetASR(ctx context.Context, audio []byte) (string, error) {
	return "mock transcription", nil
}

func NewMockASR() asr.ASRAgent {
	return &MockASR{}
}

func TestAgent(t *testing.T) {
	cfg, err := config.LoadConfig("C:/Users/BaoBao/Desktop/Project-Re-ENE/config.yaml")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	ctx := context.Background()
	llmModel := &googlegenai.GoogleAI{
		APIKey: cfg.LLMConfig.GeminiConfig.APIKey,
	}
	a := agent.NewAgent(llmModel, NewMockTTS(), NewMockASR(), nil, nil, nil, &cfg.AgentConfig, slog.Default().With("test", "TestAgent"))
	err = a.Compile(ctx)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}
	flow, err := a.InferSpeak(ctx, &agent.FlowInput{
		Text: "Tìm thông tin về sơn Tùng Mtp",
	})
	if err != nil {
		t.Fatalf("InferSpeak failed: %v", err)
	}
	for r := range flow {
		t.Log(r.Text)
	}

}
