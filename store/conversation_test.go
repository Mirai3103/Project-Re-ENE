package store

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConversationStore_CreateConversation(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	store := NewConversationStore(db)

	tests := []struct {
		name          string
		id            string
		maxWindowSize int
		characterID   string
		userID        string
		wantErr       bool
	}{
		{
			name:          "valid conversation",
			id:            uuid.New().String(),
			maxWindowSize: 10,
			characterID:   "char1",
			userID:        "user1",
			wantErr:       false,
		},
		{
			name:          "duplicate id should fail",
			id:            "duplicate-id",
			maxWindowSize: 5,
			characterID:   "char1",
			userID:        "user1",
			wantErr:       false,
		},
	}

	// Create first conversation with duplicate-id
	err := store.CreateConversation("duplicate-id", 5, "char1", "user1")
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.CreateConversation(tt.id, tt.maxWindowSize, tt.characterID, tt.userID)
			if tt.name == "duplicate id should fail" {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConversationStore_GetConversation(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	store := NewConversationStore(db)

	// Setup: Create a conversation
	conversationID := uuid.New().String()
	err := store.CreateConversation(conversationID, 10, "char1", "user1")
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{
			name:    "existing conversation",
			id:      conversationID,
			wantErr: nil,
		},
		{
			name:    "non-existing conversation",
			id:      "non-existing-id",
			wantErr: ErrNoRecordFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conv, err := store.GetConversation(tt.id)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, conv)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, conv)
				assert.Equal(t, tt.id, conv.ID)
				assert.Equal(t, 10, conv.MaxWindowSize)
				assert.Equal(t, "char1", conv.CharacterID)
				assert.Equal(t, "user1", conv.UserID)
			}
		})
	}
}

func TestConversationStore_AppendMessage(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	store := NewConversationStore(db)

	// Setup: Create a conversation
	conversationID := uuid.New().String()
	err := store.CreateConversation(conversationID, 10, "char1", "user1")
	require.NoError(t, err)

	tests := []struct {
		name           string
		conversationID string
		message        *ConversationMessage
		wantErr        bool
	}{
		{
			name:           "valid message",
			conversationID: conversationID,
			message: &ConversationMessage{
				ID:        uuid.New().String(),
				Role:      "user",
				Content:   "Hello!",
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name:           "another valid message",
			conversationID: conversationID,
			message: &ConversationMessage{
				ID:        uuid.New().String(),
				Role:      "assistant",
				Content:   "Hi there!",
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.AppendMessage(tt.conversationID, tt.message)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConversationStore_GetMessages(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	store := NewConversationStore(db)

	// Setup: Create a conversation with messages
	conversationID := uuid.New().String()
	err := store.CreateConversation(conversationID, 10, "char1", "user1")
	require.NoError(t, err)

	messages := []*ConversationMessage{
		{
			ID:        uuid.New().String(),
			Role:      "user",
			Content:   "Hello!",
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			Role:      "assistant",
			Content:   "Hi there!",
			CreatedAt: time.Now().Add(1 * time.Second),
		},
		{
			ID:        uuid.New().String(),
			Role:      "user",
			Content:   "How are you?",
			CreatedAt: time.Now().Add(2 * time.Second),
		},
	}

	for _, msg := range messages {
		err := store.AppendMessage(conversationID, msg)
		require.NoError(t, err)
	}

	tests := []struct {
		name           string
		conversationID string
		wantCount      int
		wantErr        error
	}{
		{
			name:           "existing conversation with messages",
			conversationID: conversationID,
			wantCount:      3,
			wantErr:        nil,
		},
		{
			name:           "non-existing conversation",
			conversationID: "non-existing-id",
			wantCount:      0,
			wantErr:        nil, // Empty result is not an error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgs, err := store.GetMessages(tt.conversationID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Len(t, msgs, tt.wantCount)
			}
		})
	}
}

func TestConversationStore_GetLimitedMessages(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	store := NewConversationStore(db)

	// Setup: Create a conversation with limited window size
	conversationID := uuid.New().String()
	maxWindowSize := 2
	err := store.CreateConversation(conversationID, maxWindowSize, "char1", "user1")
	require.NoError(t, err)

	// Add 5 messages
	for i := 0; i < 5; i++ {
		msg := &ConversationMessage{
			ID:        uuid.New().String(),
			Role:      "user",
			Content:   "Message " + string(rune(i)),
			CreatedAt: time.Now().Add(time.Duration(i) * time.Second),
		}
		err := store.AppendMessage(conversationID, msg)
		require.NoError(t, err)
	}

	tests := []struct {
		name           string
		conversationID string
		wantCount      int
		wantErr        error
	}{
		{
			name:           "limited messages by window size",
			conversationID: conversationID,
			wantCount:      maxWindowSize,
			wantErr:        nil,
		},
		{
			name:           "non-existing conversation",
			conversationID: "non-existing-id",
			wantCount:      0,
			wantErr:        ErrNoRecordFound, // This should error because it queries conversation first
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgs, err := store.GetLimitedMessages(tt.conversationID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Len(t, msgs, tt.wantCount)
			}
		})
	}
}
