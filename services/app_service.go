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
	cfg           *config.Config
	logger        *slog.Logger
	audioRecorder audio.Recorder
	ag            *agent.Agent
}

func NewAppService(cfg *config.Config, logger *slog.Logger, audioRecorder audio.Recorder, ag *agent.Agent) *AppService {
	return &AppService{cfg: cfg, logger: logger, audioRecorder: audioRecorder, ag: ag}
}

func (a *AppService) InvokeWithAudio(ctx context.Context, conversationID string, audioPath string) error {
	if a.audioRecorder.IsRecording() {
		return errors.New("au recorder is recording")
	}
	au, err := os.ReadFile(audioPath)
	if err != nil {
		return err
	}
	speakChan, err := a.ag.InferSpeak(ctx, &agent.FlowInput{
		Audio:          au,
		CharacterID:    "1",
		UserID:         "huuhoang",
		ConversationID: conversationID,
	})
	if err != nil {
		return err
	}
	err = a.processStreamingResponses(ctx, speakChan, func() {
		a.logger.Info("Streaming responses completed")
	})
	if err != nil {
		a.logger.Error("Error processing streaming responses", "error", err)
	}
	return nil
}
func (a *AppService) InvokeWithText(ctx context.Context, conversationID string, text string) error {
	speakChan, err := a.ag.InferSpeak(ctx, &agent.FlowInput{
		Text:           text,
		CharacterID:    "1",
		UserID:         "huuhoang",
		ConversationID: conversationID,
	})
	if err != nil {
		return err
	}
	err = a.processStreamingResponses(ctx, speakChan, func() {
		a.logger.Info("Streaming responses completed")
	})
	if err != nil {
		a.logger.Error("Error processing streaming responses", "error", err)
	}
	return nil
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
	application.Get().Event.Emit("live2d:play-audio", PlayAudioData{
		Text:   "",
		Base64: "",
		IsDone: true,
	})
	return nil
}

type PlayAudioData struct {
	Text   string
	Base64 string
	IsDone bool
}
