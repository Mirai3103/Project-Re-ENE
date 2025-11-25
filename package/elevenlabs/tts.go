package elevenlabs

import (
	"context"
	"fmt"
	"io"
)

type VoiceSettings struct {
	Stability       *float64 `json:"stability,omitempty"`
	SimilarityBoost *float64 `json:"similarity_boost,omitempty"`
	Style           *float64 `json:"style,omitempty"`
	UseSpeakerBoost *bool    `json:"use_speaker_boost,omitempty"`
	Speed           *float64 `json:"speed,omitempty"`
}

type PronunciationDictionaryLocator struct {
	PronunciationDictionaryID string  `json:"pronunciation_dictionary_id"`
	VersionID                 *string `json:"version_id,omitempty"`
}

type TTSOptions struct {
	Text                            string                            `json:"text"` // Required
	ModelID                         *string                           `json:"model_id,omitempty"`
	LanguageCode                    *string                           `json:"language_code,omitempty"`
	VoiceSettings                   *VoiceSettings                    `json:"voice_settings,omitempty"`
	PronunciationDictionaryLocators *[]PronunciationDictionaryLocator `json:"pronunciation_dictionary_locators,omitempty"`
	Seed                            *int                              `json:"seed,omitempty"`
	PreviousText                    *string                           `json:"previous_text,omitempty"`
	NextText                        *string                           `json:"next_text,omitempty"`
	PreviousRequestIDs              *[]string                         `json:"previous_request_ids,omitempty"`
	NextRequestIDs                  *[]string                         `json:"next_request_ids,omitempty"`
	ApplyTextNormalization          *string                           `json:"apply_text_normalization,omitempty"` // "auto" | "on" | "off"
	ApplyLanguageTextNormalization  *bool                             `json:"apply_language_text_normalization,omitempty"`
	UsePVCAsIVC                     *bool                             `json:"use_pvc_as_ivc,omitempty"`
	VoiceID                         string                            `json:"voice_id"`
	OutputFormat                    OutputFormat                      `json:"output_format"`
}
type OutputFormat string

const (
	// MP3 formats
	OutputFormatMP3_22050_32  OutputFormat = "mp3_22050_32"
	OutputFormatMP3_24000_48  OutputFormat = "mp3_24000_48"
	OutputFormatMP3_44100_32  OutputFormat = "mp3_44100_32"
	OutputFormatMP3_44100_64  OutputFormat = "mp3_44100_64"
	OutputFormatMP3_44100_96  OutputFormat = "mp3_44100_96"
	OutputFormatMP3_44100_128 OutputFormat = "mp3_44100_128"
	OutputFormatMP3_44100_192 OutputFormat = "mp3_44100_192"

	// PCM formats
	OutputFormatPCM_8000  OutputFormat = "pcm_8000"
	OutputFormatPCM_16000 OutputFormat = "pcm_16000"
	OutputFormatPCM_22050 OutputFormat = "pcm_22050"
	OutputFormatPCM_24000 OutputFormat = "pcm_24000"
	OutputFormatPCM_32000 OutputFormat = "pcm_32000"
	OutputFormatPCM_44100 OutputFormat = "pcm_44100"
	OutputFormatPCM_48000 OutputFormat = "pcm_48000"

	// μ-law (u-law) and a-law formats
	OutputFormatULaw_8000 OutputFormat = "ulaw_8000"
	OutputFormatALaw_8000 OutputFormat = "alaw_8000"

	// Opus formats
	OutputFormatOpus_48000_32  OutputFormat = "opus_48000_32"
	OutputFormatOpus_48000_64  OutputFormat = "opus_48000_64"
	OutputFormatOpus_48000_96  OutputFormat = "opus_48000_96"
	OutputFormatOpus_48000_128 OutputFormat = "opus_48000_128"
	OutputFormatOpus_48000_192 OutputFormat = "opus_48000_192"
)

func (c *Client) TTS(ctx context.Context, options TTSOptions) (io.Reader, error) {
	log := c.logger
	log.Debug("Generating TTS")
	resp, err := c.req.R().
		SetBody(options).
		SetContext(ctx).
		Post(fmt.Sprintf("/text-to-speech/%s?output_format=%s", options.VoiceID, options.OutputFormat))
	if err != nil {
		log.Error("Failed to generate TTS", "err", err)
		return nil, err
	}
	return resp.Body, nil
}

func (c *Client) TTSStream(ctx context.Context, options TTSOptions) (io.Reader, error) {
	log := c.logger
	log.Debug("Generating TTS stream")
	resp, err := c.req.R().
		SetBody(options).
		SetContext(ctx).
		Post(fmt.Sprintf("/text-to-speech/%s/stream?output_format=%s", options.VoiceID, options.OutputFormat))
	if err != nil {
		log.Error("Failed to generate TTS stream", "err", err)
		return nil, err
	}
	return resp.Body, nil
}
