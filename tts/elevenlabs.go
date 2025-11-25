package tts

import (
	"context"
	"io"
	"log/slog"

	ttsConfig "github.com/Mirai3103/Project-Re-ENE/config/tts"
	"github.com/Mirai3103/Project-Re-ENE/package/elevenlabs"
	"github.com/Mirai3103/Project-Re-ENE/package/utils"
)

type elevenlabsTTSAgent struct {
	client          *elevenlabs.Client
	cfg             *ttsConfig.ElevenLabsConfig
	cachingTTSAgent CachingTTSAgent
	logger          *slog.Logger
}

func newElevenlabsTTSAgent(cfg *ttsConfig.ElevenLabsConfig, cachingTTSAgent CachingTTSAgent, logger *slog.Logger) TTSAgent {
	client := elevenlabs.NewClient(elevenlabs.NewClientOptions{
		APIKey: cfg.APIKey,
	}, logger)
	return &elevenlabsTTSAgent{client: client, cfg: cfg, cachingTTSAgent: cachingTTSAgent, logger: logger}
}
func (a *elevenlabsTTSAgent) GetTTS(ctx context.Context, text string) ([]byte, error) {
	log := a.logger
	log.Debug("Getting TTS")
	cacheKey := (text + a.cfg.VoiceID + a.cfg.ModelID)
	audioBuffer := a.cachingTTSAgent.GetCachedAudioBuffer(cacheKey)
	if audioBuffer != nil {
		return audioBuffer, nil
	}
	reader, err := a.client.TTS(ctx, elevenlabs.TTSOptions{
		Text:         text,
		VoiceID:      a.cfg.VoiceID,
		ModelID:      utils.Ptr(a.cfg.ModelID),
		OutputFormat: elevenlabs.OutputFormatMP3_44100_128,
		LanguageCode: utils.Ptr("vi"),
	})
	if err != nil {
		log.Error("Failed to get TTS", "err", err)
		return nil, err
	}
	audioBuffer, err = io.ReadAll(reader)
	if err != nil {
		log.Error("Failed to read TTS", "err", err)
		return nil, err
	}
	_ = a.cachingTTSAgent.SaveCachedAudioBuffer(cacheKey, audioBuffer)

	return audioBuffer, nil
}
