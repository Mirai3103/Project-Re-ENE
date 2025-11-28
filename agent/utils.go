package agent

import (
	"context"
	"encoding/json"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
)

func (a *Agent) SaveConversationMiddleware(next core.StreamingFunc[*ai.ModelRequest, *ai.ModelResponse, *ai.ModelResponseChunk]) core.StreamingFunc[*ai.ModelRequest, *ai.ModelResponse, *ai.ModelResponseChunk] {
	return func(ctx context.Context, req *ai.ModelRequest, cb core.StreamCallback[*ai.ModelResponseChunk]) (*ai.ModelResponse, error) {
		// Trước khi chạy
		jsonData, _ := json.Marshal(req.Messages)
		a.logger.Info(" Trước khi chạy", "json", string(jsonData))
		a.logger.Info("====================================================")

		// Gọi hàm gốc
		resp, err := next(ctx, req, cb)
		if err != nil {
			a.logger.Error("Lỗi khi gọi hàm gốc", "error", err)
			return nil, err
		}
		// Sau khi chạy
		jsonData, _ = json.Marshal(resp.Message)
		a.logger.Info("Sau khi chạy", "json", string(jsonData), "error", err)
		a.logger.Info("====================================================")

		return resp, err
	}
}
