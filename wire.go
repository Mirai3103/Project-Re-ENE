//go:build wireinject

package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Mirai3103/Project-Re-ENE/agent"
	"github.com/Mirai3103/Project-Re-ENE/asr"
	"github.com/Mirai3103/Project-Re-ENE/config"
	"github.com/Mirai3103/Project-Re-ENE/llm"
	"github.com/Mirai3103/Project-Re-ENE/package/audio"
	"github.com/Mirai3103/Project-Re-ENE/services"
	"github.com/Mirai3103/Project-Re-ENE/store"
	"github.com/Mirai3103/Project-Re-ENE/tts"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/google/wire"
	"github.com/lmittmann/tint"
)

// ProvideLogger creates a new logger instance
func ProvideLogger() *slog.Logger {
	w := os.Stderr
	return slog.New(tint.NewHandler(w, nil))
}

// ProvideAudioRecorder creates a new audio recorder
func ProvideAudioRecorder(cfg *config.Config) (audio.Recorder, error) {
	return audio.NewFFmpegRecorder(audio.RecorderConfig{
		Channels:    1,
		SampleRate:  44100,
		InputDevice: cfg.ASRConfig.InputDevice,
	})
}

// ProvideLLMModel wraps the llm.New function
func ProvideLLMModel(ctx context.Context, cfg *config.Config) (*genkit.Genkit, error) {
	model, _, err := llm.New(ctx, cfg)
	return model, err
}

// ProvideLLMModelArg provides the model argument
func ProvideLLMModelArg(ctx context.Context, cfg *config.Config) (ai.ModelArg, error) {
	_, modelArg, err := llm.New(ctx, cfg)
	return modelArg, err
}

// ProvideAgentConfig extracts agent config from main config
func ProvideAgentConfig(cfg *config.Config) *config.AgentConfig {
	return &cfg.AgentConfig
}

// Application holds all initialized services
type Application struct {
	AppService      *services.AppService
	ModelService    *services.ModelService
	RecorderService *services.RecorderService
	ConfigService   *services.ConfigService
	ChatService     *services.ChatService
	Agent           *agent.Agent
}

// InitializeApplication wires up all dependencies
func InitializeApplication(ctx context.Context, cfg *config.Config) (*Application, error) {
	wire.Build(
		// Infrastructure
		ProvideLogger,
		store.NewDB,
		store.New,
		ProvideAudioRecorder,

		// Store components
		wire.FieldsOf(new(*store.Store), "ConversationStore", "CharacterStore", "UserStore"),

		// Config
		ProvideAgentConfig,

		// Agents and Models
		asr.New,
		tts.New,
		ProvideLLMModel,
		ProvideLLMModelArg,
		agent.NewAgent,

		// Services
		services.NewAppService,
		services.NewModelService,
		services.NewRecorderService,
		services.NewConfigService,
		services.NewChatService,

		// Application
		wire.Struct(new(Application), "*"),
	)
	return nil, nil
}
