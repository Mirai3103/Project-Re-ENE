package store

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	store := New(db)

	t.Run("store is initialized correctly", func(t *testing.T) {
		assert.NotNil(t, store)
		assert.NotNil(t, store.ConversationStore)
		assert.NotNil(t, store.CharacterStore)
		assert.NotNil(t, store.UserStore)
	})

	t.Run("all stores are functional", func(t *testing.T) {
		// Test ConversationStore
		err := store.CreateConversation("test-conv-1", 10, "char1", "user1")
		assert.NoError(t, err)

		conv, err := store.GetConversation("test-conv-1")
		assert.NoError(t, err)
		assert.Equal(t, "test-conv-1", conv.ID)

		// Test CharacterStore (using seeded data from migrations)
		chars, err := store.GetCharacters()
		assert.NoError(t, err)
		assert.NotEmpty(t, chars)

		// Test UserStore (using seeded data from migrations)
		// The migration 004.sql inserts a user with id 'huuhoang'
		user, err := store.GetUser("huuhoang")
		assert.NoError(t, err)
		assert.Equal(t, "huuhoang", user.ID)
	})
}

func TestCleanError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr error
	}{
		{
			name:    "sql.ErrNoRows should return ErrNoRecordFound",
			err:     sql.ErrNoRows,
			wantErr: ErrNoRecordFound,
		},
		{
			name:    "other errors should be returned as-is",
			err:     assert.AnError,
			wantErr: assert.AnError,
		},
		{
			name:    "nil error should return nil",
			err:     nil,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cleanError(tt.err)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else if tt.wantErr == ErrNoRecordFound {
				assert.ErrorIs(t, err, ErrNoRecordFound)
			} else {
				assert.Equal(t, tt.wantErr, err)
			}
		})
	}
}

func TestNewDB(t *testing.T) {
	// This test creates an actual database file
	// We should be careful with this in CI/CD environments
	t.Run("database initialization", func(t *testing.T) {
		db, err := NewDB()
		require.NoError(t, err)
		require.NotNil(t, db)
		defer db.Close()

		// Verify that tables exist by querying them
		var count int
		err = db.Get(&count, "SELECT COUNT(*) FROM conversations")
		assert.NoError(t, err)

		err = db.Get(&count, "SELECT COUNT(*) FROM characters")
		assert.NoError(t, err)

		err = db.Get(&count, "SELECT COUNT(*) FROM users")
		assert.NoError(t, err)
	})
}
