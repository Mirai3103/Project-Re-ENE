package asr

import (
	"context"
	"log/slog"

	"github.com/Mirai3103/Project-Re-ENE/config/asr"
	"github.com/Mirai3103/Project-Re-ENE/package/elevenlabs"
	"github.com/Mirai3103/Project-Re-ENE/package/utils"
)

type elevenlabsASRAgent struct {
	client *elevenlabs.Client
	cfg    *asr.ElevenLabsConfig
	logger *slog.Logger
}

func newElevenlabsASRAgent(cfg *asr.ElevenLabsConfig, logger *slog.Logger) ASRAgent {
	client := elevenlabs.NewClient(elevenlabs.NewClientOptions{
		APIKey: cfg.APIKey,
	}, logger)
	return &elevenlabsASRAgent{client: client, cfg: cfg, logger: logger}
}
func (a *elevenlabsASRAgent) GetASR(ctx context.Context, audioData []byte) (string, error) {

	response, err := a.client.CreateTranscript(ctx, audioData, elevenlabs.CreateTranscriptOptions{
		ModelID:        a.cfg.ModelID,
		LanguageCode:   utils.Ptr(a.cfg.LanguageCode),
		TagAudioEvents: utils.Ptr(false),
	})
	if err != nil {
		return "", err
	}
	return response.Text, nil
}
