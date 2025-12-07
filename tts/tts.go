package tts

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirai3103/Project-Re-ENE/config"
)

type TTSAgent interface {
	GetTTS(ctx context.Context, text string) ([]byte, error)
}

type ttsProvider struct {
	cfg             *config.Config
	cachingTTSAgent CachingTTSAgent
	logger          *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) (TTSAgent, error) {
	cachingTTSAgent := NewHashCachingTTSAgent(".cache/tts")
	switch cfg.TTSConfig.Provider {
	case "elevenlabs":
		return newElevenlabsTTSAgent(cfg.TTSConfig.ElevenLabsConfig, cachingTTSAgent, logger), nil
	default:
		return nil, fmt.Errorf("tts provider not found")
	}
}
