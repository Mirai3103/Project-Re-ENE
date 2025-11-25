package audio_test

import (
	"os"
	"testing"
	"time"

	"github.com/Mirai3103/Project-Re-ENE/package/audio"
)

func TestFFmpegRecorder_Recording4Seconds(t *testing.T) {
	rec, err := audio.NewFFmpegRecorder(audio.RecorderConfig{
		SampleRate:  44100,
		Channels:    1,
		InputDevice: "Microphone (High Definition Audio Device)",
	})
	if err != nil {
		t.Fatalf("failed to init recorder: %v", err)
	}

	t.Log("📢 Bắt đầu ghi âm, hãy nói vào micro trong 4 giây...")

	if err := rec.Start(); err != nil {
		t.Fatalf("failed to start recorder: %v", err)
	}

	time.Sleep(4 * time.Second)

	if err := rec.Stop(); err != nil {
		t.Fatalf("failed to stop recorder: %v", err)
	}

	data := rec.GetLatestAudio()
	if len(data) < 2000 {
		t.Fatalf("audio too small → có thể FFmpeg chưa ghi hoặc thiết bị sai. size=%d", len(data))
	}

	// lưu file để nghe lại
	if err := os.WriteFile("test_output.wav", data, 0644); err != nil {
		t.Fatalf("failed to write wav: %v", err)
	}

	t.Logf("✅ Đã ghi âm xong: %d bytes → test_output.wav", len(data))
}
