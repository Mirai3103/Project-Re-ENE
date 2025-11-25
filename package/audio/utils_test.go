package audio_test

import (
	"testing"

	"github.com/Mirai3103/Project-Re-ENE/package/audio"
)

func TestAudio(t *testing.T) {
	devices, err := audio.GetAudioInputDevices()
	// should at least have one device
	if err != nil {
		t.Fatalf("Failed to get audio input devices: %v", err)
	}
	for _, device := range devices {
		t.Logf("\nDevice: %s", device.Name)
	}
	if len(devices) == 0 {
		t.Fatalf("No audio input devices found")
	}
}
