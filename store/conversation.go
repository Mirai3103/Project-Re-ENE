package store

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ConversationStore struct {
	db *sqlx.DB
}

// GetConversation implements ConversationStore.
func (c *ConversationStore) GetConversation(id string) (*Conversation, error) {
	var conversation Conversation
	err := c.db.Get(&conversation, "SELECT * FROM conversations WHERE id = ?", id)
	if err != nil {
		return nil, cleanError(err)
	}
	return &conversation, nil
}

// AppendMessage implements ConversationStore.
func (c *ConversationStore) AppendMessage(conversationId string, message *ConversationMessage) error {
	_, err := c.db.Exec("INSERT INTO conversation_messages (id, conversation_id, content, role, created_at) VALUES (?, ?, ?, ?, ?)", uuid.New().String(), conversationId, message.Content, message.Role, message.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// CreateConversation implements ConversationStore.
func (c *ConversationStore) CreateConversation(Id string, maxWindowSize int, characterID string, userID string) error {
	_, err := c.db.Exec("INSERT INTO conversations (id, max_window_size, character_id, user_id) VALUES (?,?,?,?)", Id, maxWindowSize, characterID, userID)
	if err != nil {
		return err
	}
	return nil
}
func (c *ConversationStore) CreateConversationIfNotExists(Id string, maxWindowSize int, characterID string, userID string) error {
	_, err := c.GetConversation(Id)
	if errors.Is(err, ErrNoRecordFound) {
		return c.CreateConversation(Id, maxWindowSize, characterID, userID)
	}
	return err
}

func (c *ConversationStore) GetMessages(conversationId string) ([]*ConversationMessage, error) {
	rows := []*ConversationMessage{}
	err := c.db.Select(&rows, "SELECT * FROM conversation_messages WHERE conversation_id = ?", conversationId)
	if err != nil {
		return nil, cleanError(err)
	}
	return rows, nil
}

func (c *ConversationStore) GetLimitedMessages(conversationId string) ([]*ConversationMessage, error) {
	// find max_window_size from conversations table
	var maxWindowSize int
	err := c.db.Get(&maxWindowSize, "SELECT max_window_size FROM conversations WHERE id = ? LIMIT 1", conversationId)
	if err != nil {
		return nil, cleanError(err)
	}

	rows := []*ConversationMessage{}
	err = c.db.Select(&rows, "SELECT * FROM conversation_messages WHERE conversation_id = ? ORDER BY created_at DESC LIMIT ?", conversationId, maxWindowSize)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func NewConversationStore(db *sqlx.DB) *ConversationStore {
	return &ConversationStore{db: db}
}
