package services

import (
	"log/slog"

	"github.com/Mirai3103/Project-Re-ENE/config"
	"github.com/Mirai3103/Project-Re-ENE/store"
)

type ChatService struct {
	cfg               *config.Config
	logger            *slog.Logger
	conversationStore *store.ConversationStore
}

func NewChatService(cfg *config.Config, logger *slog.Logger, conversationStore *store.ConversationStore) *ChatService {
	return &ChatService{cfg: cfg, logger: logger, conversationStore: conversationStore}
}
func (s *ChatService) GetChatHistory(conversationID string) ([]*store.ConversationMessage, error) {
	return s.conversationStore.GetMessages(conversationID)
}
