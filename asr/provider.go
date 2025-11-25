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

type asrProvider struct {
	cfg    *config.Config
	logger *slog.Logger
}

func NewASRProvider(cfg *config.Config, logger *slog.Logger) ASRProvider {
	return &asrProvider{cfg: cfg, logger: logger}
}

func (p *asrProvider) GetASRAgent() (ASRAgent, error) {
	switch p.cfg.ASRConfig.Provider {
	case "elevenlabs":
		return newElevenlabsASRAgent(p.cfg.ASRConfig.ElevenLabsConfig, p.logger), nil
	default:
		return nil, fmt.Errorf("asr provider not found")
	}
}
