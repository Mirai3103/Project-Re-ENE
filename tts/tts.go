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

type TTSProvider interface {
	GetTTSAgent() (TTSAgent, error)
}

type ttsProvider struct {
	cfg             *config.Config
	cachingTTSAgent CachingTTSAgent
	logger          *slog.Logger
}

func NewTTSProvider(cfg *config.Config, logger *slog.Logger) TTSProvider {
	cachingTTSAgent := NewHashCachingTTSAgent(".cache/tts")
	return &ttsProvider{cfg: cfg, cachingTTSAgent: cachingTTSAgent}
}

func (p *ttsProvider) GetTTSAgent() (TTSAgent, error) {
	switch p.cfg.TTSConfig.Provider {
	case "elevenlabs":
		return newElevenlabsTTSAgent(p.cfg.TTSConfig.ElevenLabsConfig, p.cachingTTSAgent, p.logger), nil
	default:
		return nil, fmt.Errorf("tts provider not found")
	}
}
