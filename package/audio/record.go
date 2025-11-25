package audio

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

var (
	ErrAlreadyRecording = errors.New("recorder already running")
	ErrNotRecording     = errors.New("recorder not running")
	ErrFFmpegNotFound   = errors.New("ffmpeg not found in PATH")
)

// Recorder defines the contract for audio recording.
type Recorder interface {

	// Start begins audio recording. It is non-blocking and returns immediately.
	Start() error

	// Stop ends audio recording.
	Stop() error

	// GetLatestAudio retrieves the most recent recorded audio data in WAV format.
	GetLatestAudio() []byte // read latest wav

	GetLatestAudioPath() string

	// IsRecording indicates whether the recorder is currently active.
	IsRecording() bool

	// ClearData removes all recorded audio data from memory.
	ClearData()

	// Close releases any resources held by the recorder.
	Close() error

	GetAvailableInputDevices() ([]Device, error)
}

// RecorderConfig holds the configuration for audio recording.
type RecorderConfig struct {
	SampleRate  uint32
	Channels    uint32
	InputDevice string
}

type ffmpegRecorder struct {
	cfg      RecorderConfig
	cmd      *exec.Cmd
	filePath string

	mu        sync.RWMutex
	recording bool
	stdin     io.WriteCloser
}

func NewFFmpegRecorder(cfg RecorderConfig) (Recorder, error) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return nil, ErrFFmpegNotFound
	}

	r := &ffmpegRecorder{
		cfg: cfg,
	}

	return r, nil
}

func (r *ffmpegRecorder) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.recording {
		return ErrAlreadyRecording
	}

	// tạo file WAV tạm
	tmpFile, err := os.CreateTemp("", "record_*.wav")
	if err != nil {
		return err
	}
	_ = tmpFile.Close()
	r.filePath = tmpFile.Name()

	// build ffmpeg command
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command(
			"ffmpeg",
			"-y", // always overwrite
			"-f", "dshow",
			"-i", `audio=`+r.cfg.InputDevice,
			"-ac", "1",
			"-ar", toStr(r.cfg.SampleRate),
			"-vn",
			r.filePath,
		)
	} else {
		cmd = exec.Command(
			"ffmpeg",
			"-y",
			"-f", "pulse",
			"-i", r.cfg.InputDevice,
			"-ac", "1",
			"-ar", toStr(r.cfg.SampleRate),
			"-vn",
			r.filePath,
		)
	}
	r.cmd = cmd
	fmt.Printf("FFmpeg command: %s\n", r.cmd.String())
	// tránh block stderr
	r.cmd.Stderr = io.Discard
	stdin, err := r.cmd.StdinPipe()
	if err != nil {
		return err
	}
	r.stdin = stdin
	if err := r.cmd.Start(); err != nil {
		return err
	}

	r.recording = true
	return nil
}

func (r *ffmpegRecorder) GetLatestAudioPath() string {
	return r.filePath
}

func (r *ffmpegRecorder) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.recording {
		return ErrNotRecording
	}

	// gửi "q" để ffmpeg tự thoát
	if r.stdin != nil {
		_, _ = r.stdin.Write([]byte("q\n"))
		defer r.stdin.Close()
	}

	done := make(chan error, 1)
	go func() {
		done <- r.cmd.Wait()
	}()

	select {
	case <-time.After(2 * time.Second):
		// ffmpeg bị treo → kill
		_ = r.cmd.Process.Kill()
		<-done
	case <-done:
		// ffmpeg tự thoát ok
	}

	r.recording = false
	return nil
}

func (r *ffmpegRecorder) GetLatestAudio() []byte {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.filePath == "" {
		return nil
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil
	}

	// tránh trả file rỗng khi FFmpeg đang ghi
	return bytes.Clone(data)
}

func (r *ffmpegRecorder) IsRecording() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.recording
}

func (r *ffmpegRecorder) ClearData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.filePath != "" {
		_ = os.Remove(r.filePath)
	}

	r.filePath = ""
}

func (r *ffmpegRecorder) Close() error {
	_ = r.Stop()
	r.ClearData()
	return nil
}

// helper
func toStr[T ~int | ~uint32](v T) string {
	return fmt.Sprintf("%d", v)
}

func (r *ffmpegRecorder) GetAvailableInputDevices() ([]Device, error) {
	return GetAudioInputDevices()
}
