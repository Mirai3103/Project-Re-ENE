package services

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/Mirai3103/Project-Re-ENE/agent"
	"github.com/Mirai3103/Project-Re-ENE/config"
	"github.com/Mirai3103/Project-Re-ENE/package/audio"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type AppService struct {
	cfg            *config.Config
	logger         *slog.Logger
	assistantAgent agent.AssistantAgent
	audioRecorder  audio.Recorder
}

func NewAppService(cfg *config.Config, logger *slog.Logger, assistantAgent agent.AssistantAgent, audioRecorder audio.Recorder) *AppService {
	return &AppService{cfg: cfg, logger: logger, assistantAgent: assistantAgent, audioRecorder: audioRecorder}
}

func (a *AppService) InvokeWithAudio(ctx context.Context, conversationID string, audioPath string) error {
	if a.audioRecorder.IsRecording() {
		return errors.New("audio recorder is recording")
	}
	audio, err := os.ReadFile(audioPath)
	if err != nil {
		return err
	}
	speakResponseStream, err := a.assistantAgent.Stream(ctx, conversationID, agent.UserInput{Audio: &audio})
	if err != nil {
		return err
	}
	return a.processStreamingResponses(ctx, speakResponseStream, func() {
		a.audioRecorder.ClearData()
	})
}
func (a *AppService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	if err := a.assistantAgent.CompileChain(ctx); err != nil {
		return err
	}
	return nil
}

type PlayAudioData struct {
	Text   string
	Base64 string
}

func (a *AppService) processStreamingResponses(ctx context.Context, stream chan agent.SpeakResponse, onDone func()) error {
	defer onDone()
	for speakResponse := range stream {
		a.logger.Info("Received speak response", "text", speakResponse.Text)

		// Emit event to frontend
		application.Get().Event.Emit("live2d:play-audio", PlayAudioData{
			Text:   speakResponse.Text,
			Base64: speakResponse.ToBase64(),
		})
	}
	onDone()
	return nil
}
