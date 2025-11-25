package audio

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

// Device AudioDevice represents an audio input device.
type Device struct {
	ID        string
	Name      string
	IsDefault bool
	Format    string // "pulse", "avfoundation", "dshow", etc.
}

// String returns a human-readable representation of the device.
func (d Device) String() string {
	defaultMarker := ""
	if d.IsDefault {
		defaultMarker = " [DEFAULT]"
	}
	return fmt.Sprintf("%s: %s%s", d.ID, d.Name, defaultMarker)
}

// GetAudioInputDevices returns a list of available audio input devices.
func GetAudioInputDevices() ([]Device, error) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return nil, ErrFFmpegNotFound
	}

	switch runtime.GOOS {
	case "linux":
		return getLinuxDevices()
	case "darwin":
		return getMacDevices()
	case "windows":
		return getWindowsDevices()
	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

// GetDefaultAudioDevice returns the default audio input device.
func GetDefaultAudioDevice() (*Device, error) {
	devices, err := GetAudioInputDevices()
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		if device.IsDefault {
			return &device, nil
		}
	}

	if len(devices) > 0 {
		return &devices[0], nil
	}

	return nil, errors.New("no audio input devices found")
}

// getLinuxDevices lists audio devices on Linux (PulseAudio/ALSA).
func getLinuxDevices() ([]Device, error) {
	// Try PulseAudio first
	devices, err := getPulseAudioDevices()
	if err == nil && len(devices) > 0 {
		return devices, nil
	}

	// Fallback to ALSA
	return getALSADevices()
}

// getPulseAudioDevices lists PulseAudio devices.
func getPulseAudioDevices() ([]Device, error) {
	cmd := exec.Command("pactl", "list", "sources", "short")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var devices []Device
	scanner := bufio.NewScanner(bytes.NewReader(output))
	isFirst := true

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		// Skip monitor sources (they capture output, not input)
		if strings.Contains(fields[1], ".monitor") {
			continue
		}

		device := Device{
			ID:        fields[1],
			Name:      getDeviceFriendlyName(fields[1]),
			IsDefault: isFirst,
			Format:    "pulse",
		}
		devices = append(devices, device)
		isFirst = false
	}

	if len(devices) == 0 {
		// Add default device as fallback
		devices = append(devices, Device{
			ID:        "default",
			Name:      "Default PulseAudio Device",
			IsDefault: true,
			Format:    "pulse",
		})
	}

	return devices, nil
}

// getALSADevices lists ALSA devices.
func getALSADevices() ([]Device, error) {
	cmd := exec.Command("arecord", "-L")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var devices []Device
	scanner := bufio.NewScanner(bytes.NewReader(output))
	isFirst := true

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, " ") {
			continue
		}

		// Skip null devices
		if strings.HasPrefix(line, "null") {
			continue
		}

		device := Device{
			ID:        line,
			Name:      getDeviceFriendlyName(line),
			IsDefault: isFirst && line == "default",
			Format:    "alsa",
		}
		devices = append(devices, device)
		isFirst = false
	}

	if len(devices) == 0 {
		devices = append(devices, Device{
			ID:        "default",
			Name:      "Default ALSA Device",
			IsDefault: true,
			Format:    "alsa",
		})
	}

	return devices, nil
}

// getMacDevices lists audio devices on macOS (AVFoundation).
func getMacDevices() ([]Device, error) {
	cmd := exec.Command("ffmpeg", "-f", "avfoundation", "-list_devices", "true", "-i", "")
	output, err := cmd.CombinedOutput()
	if err == nil {
		return nil, errors.New("ffmpeg list devices should return error")
	}

	var devices []Device
	scanner := bufio.NewScanner(bytes.NewReader(output))

	// Regex to match: [AVFoundation indev @ 0x...] [0] Built-in Microphone
	deviceRegex := regexp.MustCompile(`\[AVFoundation.*?\]\s+\[(\d+)\]\s+(.+)`)
	inAudioSection := false

	for scanner.Scan() {
		line := scanner.Text()

		// Check if we're in the audio input section
		if strings.Contains(line, "AVFoundation audio devices:") {
			inAudioSection = true
			continue
		}

		// Stop when we hit video devices
		if strings.Contains(line, "AVFoundation video devices:") {
			break
		}

		if !inAudioSection {
			continue
		}

		matches := deviceRegex.FindStringSubmatch(line)
		if len(matches) == 3 {
			deviceID := matches[1]
			deviceName := strings.TrimSpace(matches[2])

			device := Device{
				ID:        ":" + deviceID,
				Name:      deviceName,
				IsDefault: deviceID == "0",
				Format:    "avfoundation",
			}
			devices = append(devices, device)
		}
	}

	if len(devices) == 0 {
		devices = append(devices, Device{
			ID:        ":0",
			Name:      "Default Audio Device",
			IsDefault: true,
			Format:    "avfoundation",
		})
	}

	return devices, nil
}

// getWindowsDevices lists audio devices on Windows (DirectShow).
func getWindowsDevices() ([]Device, error) {
	cmd := exec.Command("ffmpeg", "-f", "dshow", "-list_devices", "true", "-i", "dummy")
	output, _ := cmd.CombinedOutput()

	var devices []Device
	scanner := bufio.NewScanner(bytes.NewReader(output))

	// match --> "xxx" (audio)
	deviceRegex := regexp.MustCompile(`"([^"]+)"\s+\(audio\)`)

	for scanner.Scan() {
		line := scanner.Text()

		// Chỉ cần tìm dòng chứa (audio)
		if !strings.Contains(line, "(audio)") {
			continue
		}

		matches := deviceRegex.FindStringSubmatch(line)
		if len(matches) == 2 {
			name := strings.TrimSpace(matches[1])

			devices = append(devices, Device{
				ID:        "audio=" + name,
				Name:      name,
				IsDefault: len(devices) == 0,
				Format:    "dshow",
			})
		}
	}

	if len(devices) == 0 {
		devices = append(devices, Device{
			ID:        "audio=Microphone",
			Name:      "Default Microphone",
			IsDefault: true,
			Format:    "dshow",
		})
	}

	return devices, nil
}

// getDeviceFriendlyName extracts a friendly name from device ID.
func getDeviceFriendlyName(deviceID string) string {
	// Remove common prefixes
	name := deviceID
	name = strings.TrimPrefix(name, "alsa_input.")
	name = strings.TrimPrefix(name, "alsa_output.")

	// Replace underscores and dots with spaces
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, ".", " ")

	// Capitalize first letter
	if len(name) > 0 {
		name = strings.ToUpper(string(name[0])) + name[1:]
	}

	return name
}
