package agent

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/Mirai3103/Project-Re-ENE/package/utils"
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
		jsonData, _ := json.Marshal(lastMessage)
		err := a.store.CreateConversationMessage(ctx, store.CreateConversationMessageParams{
			ConversationID: utils.Ptr(conversationID),
			Role:           utils.Ptr(string(lastMessage.Role)),
			Content:        jsonData,
			CreatedAt:      utils.Ptr(time.Now()),
			ID:             uuid.New().String(),
		})
		if err != nil {
			a.logger.Error("Lỗi khi lưu tin nhắn", "error", err)
			return nil, err
		}

		// Gọi hàm gốc
		a.logger.Info("Calling next function")
		resp, err := next(ctx, req, cb)
		a.logger.Info("Next function called", "resp", resp)
		if err != nil {
			a.logger.Error("Lỗi khi gọi hàm gốc", "error", err)
			return nil, err
		}
		// Sau khi chạy
		cpyMgs, err := DeepCopyMessage(resp.Message)
		if err != nil {
			a.logger.Error("Lỗi khi sao chép tin nhắn", "error", err)
			return resp, err
		}
		NormalizeMessage(cpyMgs)
		jsonData, _ = json.Marshal(cpyMgs)
		a.logger.Info("SaveAssistantMessage")
		err = a.store.CreateConversationMessage(ctx, store.CreateConversationMessageParams{
			ConversationID: utils.Ptr(conversationID),
			Role:           utils.Ptr("assistant"),
			Content:        jsonData,
			CreatedAt:      utils.Ptr(time.Now()),
			ID:             uuid.New().String(),
		})
		if err != nil {
			a.logger.Error("Lỗi khi lưu tin nhắn", "error", err)
		}

		return resp, err
	}
}
func DeepCopyMessage(src *ai.Message) (*ai.Message, error) {
	if src == nil {
		return nil, nil
	}

	b, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}

	var dst ai.Message
	if err := json.Unmarshal(b, &dst); err != nil {
		return nil, err
	}

	return &dst, nil
}

func NormalizeMessage(msg *ai.Message) {
	if msg == nil {
		return
	}

	var result []*ai.Part
	var buffer strings.Builder

	flushBuffer := func() {
		if buffer.Len() > 0 {
			result = append(result, &ai.Part{
				Kind: ai.PartText,
				Text: buffer.String(),
			})
			buffer.Reset()
		}
	}

	for _, p := range msg.Content {
		if p == nil {
			continue
		}

		switch p.Kind {

		case ai.PartText:
			// gom text
			if p.Text != "" {
				buffer.WriteString(p.Text)
			}

		case ai.PartToolRequest, ai.PartToolResponse:
			flushBuffer()
			result = append(result, p)

		default:
			flushBuffer()
			result = append(result, p)
		}
	}

	// flush phần text cuối
	flushBuffer()

	msg.Content = result
}

func ConversationToText(conversation []*ai.Message) string {
	chBuilder := strings.Builder{}
	contentBuilder := strings.Builder{}
	for _, message := range conversation {
		for _, part := range message.Content {
			contentBuilder.WriteString(part.Text)
		}
		content := contentBuilder.String()
		if content != "" {
			chBuilder.WriteString(string(message.Role))
			chBuilder.WriteString(": ")
			chBuilder.WriteString(content)
			chBuilder.WriteString("\n")
		}
		contentBuilder.Reset()
	}
	return chBuilder.String()

}

func UserFactsToText(userFacts []store.UserFact) string {
	ufBuilder := strings.Builder{}
	for _, userFact := range userFacts {
		ufBuilder.WriteString(*userFact.Name)
		ufBuilder.WriteString(": ")
		ufBuilder.WriteString(*userFact.Value)
		ufBuilder.WriteString("\n")
	}
	return ufBuilder.String()
}

func CharacterFactsToText(characterFacts []store.CharacterFact) string {
	cfBuilder := strings.Builder{}
	for _, characterFact := range characterFacts {
		cfBuilder.WriteString(*characterFact.Name)
		cfBuilder.WriteString(": ")
		cfBuilder.WriteString(*characterFact.Value)
		cfBuilder.WriteString("\n")
	}
	return cfBuilder.String()
}

func ParseHistoryMessages(messages []store.ConversationMessage) []*ai.Message {
	historyMessages := make([]*ai.Message, len(messages))
	for i, message := range messages {
		var hm ai.Message
		err := json.Unmarshal(message.Content, &hm)
		if err != nil {
			continue
		}
		hm.Role = ai.Role(*message.Role)
		historyMessages[i] = &hm
	}
	return historyMessages
}
