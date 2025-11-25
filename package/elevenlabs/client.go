package elevenlabs

import (
	"log/slog"

	"github.com/Mirai3103/Project-Re-ENE/package/utils"

	"github.com/imroc/req/v3"
)

type Client struct {
	apiKey  string
	req     *req.Client
	baseURL string
	logger  *slog.Logger
}
type NewClientOptions struct {
	APIKey        string
	BaseURL       *string
	DevMode       *bool
	CommonHeaders *map[string]string
}

func NewClient(options NewClientOptions, logger *slog.Logger) *Client {
	client := req.C()
	if utils.OrDefault(options.DevMode, false) {
		client = client.DevMode()
	}
	client = client.SetBaseURL(utils.OrDefault(options.BaseURL, "https://api.elevenlabs.io/v1"))
	var headers map[string]string
	if options.CommonHeaders != nil {
		headers = *options.CommonHeaders
		headers["xi-api-key"] = options.APIKey
	} else {
		headers = map[string]string{
			"xi-api-key": options.APIKey,
		}
	}
	client = client.SetCommonHeaders(headers)
	return &Client{
		apiKey:  options.APIKey,
		req:     client,
		baseURL: utils.OrDefault(options.BaseURL, "https://api.elevenlabs.io/v1"),
	}
}
