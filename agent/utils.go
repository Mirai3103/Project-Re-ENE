package agent

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Mirai3103/Project-Re-ENE/store"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/google/uuid"
)

func (a *Agent) SaveConversationMiddleware(next core.StreamingFunc[*ai.ModelRequest, *ai.ModelResponse, *ai.ModelResponseChunk]) core.StreamingFunc[*ai.ModelRequest, *ai.ModelResponse, *ai.ModelResponseChunk] {
	return func(ctx context.Context, req *ai.ModelRequest, cb core.StreamCallback[*ai.ModelResponseChunk]) (*ai.ModelResponse, error) {
		conversationID := ctx.Value(ConversationID).(string)
		characterID := ctx.Value(CharacterID).(string)
		userID := ctx.Value(UserID).(string)
		a.logger.Info("SaveConversationMiddleware", "conversationID", conversationID, "characterID", characterID, "userID", userID)
		// Trước khi chạy
		lastMessage := req.Messages[len(req.Messages)-1]
		a.logger.Info("SaveToolMessage", "lastMessage", lastMessage)
		jsonData, _ := json.Marshal(lastMessage)
		err := a.conversationStore.AppendMessage(conversationID, &store.ConversationMessage{
			ConversationID: conversationID,
			Role:           string(lastMessage.Role),
			Content:        string(jsonData),
			CreatedAt:      time.Now(),
			ID:             uuid.New().String(),
		})
		if err != nil {
			a.logger.Error("Lỗi khi lưu tin nhắn", "error", err)
			return nil, err
		}

		// Gọi hàm gốc
		resp, err := next(ctx, req, cb)
		if err != nil {
			a.logger.Error("Lỗi khi gọi hàm gốc", "error", err)
			return nil, err
		}
		// Sau khi chạy
		jsonData, _ = json.Marshal(resp.Message)
		a.logger.Info("SaveAssistantMessage")
		err = a.conversationStore.AppendMessage(conversationID, &store.ConversationMessage{
			ConversationID: conversationID,
			Role:           "assistant",
			Content:        string(jsonData),
			CreatedAt:      time.Now(),
			ID:             uuid.New().String(),
		})
		if err != nil {
			a.logger.Error("Lỗi khi lưu tin nhắn", "error", err)
			return nil, err
		}
		return resp, err
	}
}
