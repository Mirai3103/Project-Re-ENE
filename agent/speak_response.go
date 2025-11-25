package agent

import (
	"encoding/base64"
)

type SpeakResponse struct {
	Text        string `json:"text"`
	AudioBuffer []byte `json:"audio_buffer"`
}

func (s *SpeakResponse) ToBase64() string {
	encoded := base64.StdEncoding.EncodeToString(s.AudioBuffer)
	return encoded
}
