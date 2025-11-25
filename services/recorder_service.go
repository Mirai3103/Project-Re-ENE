package services

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/Mirai3103/Project-Re-ENE/config"
	"github.com/Mirai3103/Project-Re-ENE/package/audio"
)

type RecorderService struct {
	cfg           *config.Config
	done          chan struct{}
	audioRecorder audio.Recorder
	logger        *slog.Logger
}

var (
	ErrAlreadyRecording = errors.New("recorder already recording")
	ErrNotRecording     = errors.New("recorder not recording")
	ErrEmptyCommand     = errors.New("empty command")
	ErrMissingConvID    = errors.New("missing conversationID")
	ErrUnsupportedOS    = errors.New("unsupported operating system")
)

func NewRecorderService(recCfg *config.Config, recorder audio.Recorder) *RecorderService {
	return &RecorderService{
		cfg:           recCfg,
		audioRecorder: recorder,
		logger:        slog.Default(),
	}
}

func (a *RecorderService) StartRecording() error {
	if a.audioRecorder.IsRecording() {
		return ErrAlreadyRecording
	}

	if err := a.audioRecorder.Start(); err != nil {
		return fmt.Errorf("failed to start recording: %w", err)
	}

	a.logger.Info("Recording started")
	return nil
}

// StopRecording stops audio recording
func (a *RecorderService) StopRecording() string {
	if !a.audioRecorder.IsRecording() {
		return ""
	}

	a.audioRecorder.Stop()
	a.logger.Info("Recording stopped")
	return a.audioRecorder.GetLatestAudioPath()
}

func (a *RecorderService) GetAvailableInputDevices() ([]audio.Device, error) {
	devices, err := a.audioRecorder.GetAvailableInputDevices()
	if err != nil {
		return nil, err
	}
	return devices, nil
}
