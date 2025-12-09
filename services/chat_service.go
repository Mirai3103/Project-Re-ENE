package services

import (
	"context"
	"log/slog"

	"github.com/Mirai3103/Project-Re-ENE/config"
	"github.com/Mirai3103/Project-Re-ENE/package/utils"
	"github.com/Mirai3103/Project-Re-ENE/store"
)

type ChatService struct {
	cfg    *config.Config
	logger *slog.Logger
	store  *store.Queries
}

func NewChatService(cfg *config.Config, logger *slog.Logger, store *store.Queries) *ChatService {
	return &ChatService{cfg: cfg, logger: logger, store: store}
}
func (s *ChatService) GetChatHistory(ctx context.Context, conversationID string) ([]store.ConversationMessage, error) {
	messages, err := s.store.ListConversationMessages(ctx, utils.Ptr(conversationID))
	if err != nil {
		return nil, err
	}
	return messages, nil
}
