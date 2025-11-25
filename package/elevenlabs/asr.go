package elevenlabs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// Transcript request options
type CreateTranscriptOptions struct {
	AudioData             []byte      `json:"-"`        // Required: audio file as byte array
	ModelID               string      `json:"model_id"` // Required: "scribe_v1" or "scribe_v1_experimental"
	LanguageCode          *string     `json:"language_code,omitempty"`
	TagAudioEvents        *bool       `json:"tag_audio_events,omitempty"`       // Default: true
	NumSpeakers           *int        `json:"num_speakers,omitempty"`           // 1-32
	TimestampsGranularity *string     `json:"timestamps_granularity,omitempty"` // "none", "word", "character"
	Diarize               *bool       `json:"diarize,omitempty"`                // Default: false
	DiarizationThreshold  *float64    `json:"diarization_threshold,omitempty"`  // 0.1-0.4
	FileFormat            *string     `json:"file_format,omitempty"`            // "pcm_s16le_16" or "other"
	Temperature           *float64    `json:"temperature,omitempty"`            // 0-2
	Seed                  *int        `json:"seed,omitempty"`                   // 0-2147483647
	UseMultiChannel       *bool       `json:"use_multi_channel,omitempty"`      // Default: false
	Webhook               *bool       `json:"webhook,omitempty"`                // Default: false
	WebhookID             *string     `json:"webhook_id,omitempty"`
	WebhookMetadata       interface{} `json:"webhook_metadata,omitempty"`
	EnableLogging         *bool       `json:"-"` // Query parameter
}

// Word timing information
type TranscriptWord struct {
	Text         string   `json:"text"`
	Start        float64  `json:"start"`
	End          float64  `json:"end"`
	Type         string   `json:"type"`
	SpeakerID    *string  `json:"speaker_id,omitempty"`
	Logprob      *float64 `json:"logprob,omitempty"`
	ChannelIndex *int     `json:"channel_index,omitempty"`
}

// Additional format response
type AdditionalFormat struct {
	Format string `json:"format"`
	Text   string `json:"text"`
}

// Single channel transcript response
type TranscriptResponse struct {
	LanguageCode        string              `json:"language_code"`
	LanguageProbability float64             `json:"language_probability"`
	Text                string              `json:"text"`
	Words               []TranscriptWord    `json:"words"`
	ChannelIndex        *int                `json:"channel_index,omitempty"`
	AdditionalFormats   *[]AdditionalFormat `json:"additional_formats,omitempty"`
	TranscriptionID     *string             `json:"transcription_id,omitempty"`
}

// Multi-channel transcript response
type MultichannelTranscriptResponse struct {
	Transcripts     []TranscriptResponse `json:"transcripts"`
	TranscriptionID *string              `json:"transcription_id,omitempty"`
}

// Webhook response (for async processing)
type WebhookTranscriptResponse struct {
	Message         string  `json:"message"`
	RequestID       string  `json:"request_id"`
	TranscriptionID *string `json:"transcription_id,omitempty"`
}

// CreateTranscript transcribes audio data and returns the transcript
func (c *Client) CreateTranscript(ctx context.Context, audioData []byte, options CreateTranscriptOptions) (*TranscriptResponse, error) {
	log := c.logger
	log.Debug("Creating transcript")
	// Validate required fields
	if len(audioData) == 0 {
		log.Error("Audio data cannot be empty")
		return nil, errors.New("audio data cannot be empty")
	}
	if options.ModelID == "" {
		log.Error("Model ID is required")
		return nil, errors.New("model_id is required")
	}

	// Build multipart form request
	req := c.req.R().
		SetFileReader("file", "audio.mp3", bytes.NewReader(audioData)).
		SetContext(ctx).
		SetFormData(map[string]string{
			"model_id": options.ModelID,
		})

	// Add optional fields
	if options.LanguageCode != nil {
		req.SetFormData(map[string]string{"language_code": *options.LanguageCode})
	}
	if options.TagAudioEvents != nil {
		req.SetFormData(map[string]string{"tag_audio_events": fmt.Sprintf("%v", *options.TagAudioEvents)})
	}
	if options.NumSpeakers != nil {
		req.SetFormData(map[string]string{"num_speakers": fmt.Sprintf("%d", *options.NumSpeakers)})
	}
	if options.TimestampsGranularity != nil {
		req.SetFormData(map[string]string{"timestamps_granularity": *options.TimestampsGranularity})
	}
	if options.Diarize != nil {
		req.SetFormData(map[string]string{"diarize": fmt.Sprintf("%v", *options.Diarize)})
	}
	if options.DiarizationThreshold != nil {
		req.SetFormData(map[string]string{"diarization_threshold": fmt.Sprintf("%f", *options.DiarizationThreshold)})
	}
	if options.FileFormat != nil {
		req.SetFormData(map[string]string{"file_format": *options.FileFormat})
	}
	if options.Temperature != nil {
		req.SetFormData(map[string]string{"temperature": fmt.Sprintf("%f", *options.Temperature)})
	}
	if options.Seed != nil {
		req.SetFormData(map[string]string{"seed": fmt.Sprintf("%d", *options.Seed)})
	}
	if options.UseMultiChannel != nil {
		req.SetFormData(map[string]string{"use_multi_channel": fmt.Sprintf("%v", *options.UseMultiChannel)})
	}
	if options.Webhook != nil {
		req.SetFormData(map[string]string{"webhook": fmt.Sprintf("%v", *options.Webhook)})
	}
	if options.WebhookID != nil {
		req.SetFormData(map[string]string{"webhook_id": *options.WebhookID})
	}
	if options.WebhookMetadata != nil {
		metadataJSON, _ := json.Marshal(options.WebhookMetadata)
		req.SetFormData(map[string]string{"webhook_metadata": string(metadataJSON)})
	}

	// Add query parameter
	if options.EnableLogging != nil {
		req.SetQueryParam("enable_logging", fmt.Sprintf("%v", *options.EnableLogging))
	}

	// Make request
	resp, err := req.Post("/speech-to-text")
	if err != nil {
		log.Error("Failed to create transcript", "err", err)
		return nil, errors.New("request failed")
	}

	if !resp.IsSuccessState() {
		log.Error("Failed to create transcript", "status", resp.GetStatusCode(), "body", resp.String())
		return nil, errors.New("API error")
	}

	// Parse response
	var result TranscriptResponse
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		log.Error("Failed to parse response", "err", err)
		return nil, errors.New("failed to parse response")
	}

	return &result, nil
}

// CreateTranscriptMultichannel is a helper for multi-channel audio
func (c *Client) CreateTranscriptMultichannel(ctx context.Context, audioData []byte, options CreateTranscriptOptions) (*MultichannelTranscriptResponse, error) {
	log := c.logger
	log.Debug("Creating transcript multichannel")
	// Force multi-channel mode
	useMultiChannel := true
	options.UseMultiChannel = &useMultiChannel

	// Similar implementation but parse as MultichannelTranscriptResponse
	req := c.req.R().
		SetFileReader("file", "audio.mp3", bytes.NewReader(audioData)).
		SetContext(ctx).
		SetFormData(map[string]string{
			"model_id":          options.ModelID,
			"use_multi_channel": "true",
		})

	// Add other optional fields (same as above)
	// ... (implement similar to CreateTranscript)

	resp, err := req.Post("/speech-to-text")
	if err != nil {
		log.Error("Failed to create transcript multichannel", "err", err)
		return nil, errors.New("request failed")
	}

	if !resp.IsSuccessState() {
		log.Error("Failed to create transcript multichannel", "status", resp.GetStatusCode(), "body", resp.String())
		return nil, errors.New("API error")
	}

	var result MultichannelTranscriptResponse
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		log.Error("Failed to parse response", "err", err)
		return nil, errors.New("failed to parse response")
	}

	return &result, nil
}
