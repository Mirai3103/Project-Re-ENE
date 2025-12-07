package asr

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirai3103/Project-Re-ENE/config"
)

type ASRAgent interface {
	GetASR(ctx context.Context, audioData []byte) (string, error)
}

type ASRProvider interface {
	GetASRAgent() (ASRAgent, error)
}

func New(cfg *config.Config, logger *slog.Logger) (ASRAgent, error) {
	switch cfg.ASRConfig.Provider {
	case "elevenlabs":
		return newElevenlabsASRAgent(cfg.ASRConfig.ElevenLabsConfig, logger), nil
	default:
		return nil, fmt.Errorf("asr provider not found")
	}
}
