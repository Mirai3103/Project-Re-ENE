package store

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupCharacterTestData(t *testing.T, db *sqlx.DB) (string, string) {
	t.Helper()

	characterID := uuid.New().String()
	characterID2 := uuid.New().String()

	// Insert test characters
	_, err := db.Exec(`
		INSERT INTO characters (id, name, base_prompt, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, characterID, "Test Character", "Test prompt", "Test description", time.Now(), time.Now())
	require.NoError(t, err)

	_, err = db.Exec(`
		INSERT INTO characters (id, name, base_prompt, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, characterID2, "Another Character", "Another prompt", "Another description", time.Now(), time.Now())
	require.NoError(t, err)

	// Insert test character facts
	for i := 0; i < 3; i++ {
		factID := uuid.New().String()
		_, err := db.Exec(`
			INSERT INTO character_facts (id, character_id, name, value, type, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, factID, characterID, "fact"+string(rune(i)), "value"+string(rune(i)), "type"+string(rune(i)),
			time.Now().Add(time.Duration(i)*time.Second), time.Now().Add(time.Duration(i)*time.Second))
		require.NoError(t, err)
	}

	return characterID, characterID2
}

func TestCharacterStore_GetCharacter(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	store := NewCharacterStore(db)
	characterID, _ := setupCharacterTestData(t, db)

	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{
			name:    "existing character",
			id:      characterID,
			wantErr: nil,
		},
		{
			name:    "non-existing character",
			id:      "non-existing-id",
			wantErr: ErrNoRecordFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char, err := store.GetCharacter(tt.id)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, char)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, char)
				assert.Equal(t, tt.id, char.ID)
				assert.Equal(t, "Test Character", char.Name)
				assert.Equal(t, "Test prompt", char.BasePrompt)
				assert.Equal(t, "Test description", char.Description)
			}
		})
	}
}

func TestCharacterStore_GetCharacters(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	store := NewCharacterStore(db)

	t.Run("empty database", func(t *testing.T) {
		chars, err := store.GetCharacters()
		assert.NoError(t, err)
		// Database might have seed data from migrations
		assert.NotNil(t, chars)
	})

	t.Run("with characters", func(t *testing.T) {
		setupCharacterTestData(t, db)
		chars, err := store.GetCharacters()
		assert.NoError(t, err)
		assert.NotNil(t, chars)
		// Should have at least the 2 we added plus any seed data
		assert.GreaterOrEqual(t, len(chars), 2)
	})
}

func TestCharacterStore_GetCharacterFacts(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	store := NewCharacterStore(db)
	characterID, _ := setupCharacterTestData(t, db)

	tests := []struct {
		name          string
		characterID   string
		limit         int
		wantMinCount  int
		wantMaxCount  int
		wantErr       error
		checkOrdering bool
	}{
		{
			name:          "get all facts",
			characterID:   characterID,
			limit:         10,
			wantMinCount:  3,
			wantMaxCount:  3,
			wantErr:       nil,
			checkOrdering: true,
		},
		{
			name:         "limit to 2 facts",
			characterID:  characterID,
			limit:        2,
			wantMinCount: 2,
			wantMaxCount: 2,
			wantErr:      nil,
		},
		{
			name:         "non-existing character",
			characterID:  "non-existing-id",
			limit:        10,
			wantMinCount: 0,
			wantMaxCount: 0,
			wantErr:      nil, // Empty result is not an error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			facts, err := store.GetCharacterFacts(tt.characterID, tt.limit)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, facts)
				assert.GreaterOrEqual(t, len(facts), tt.wantMinCount)
				assert.LessOrEqual(t, len(facts), tt.wantMaxCount)

				// Check ordering (should be DESC by updated_at)
				if tt.checkOrdering && len(facts) > 1 {
					for i := 0; i < len(facts)-1; i++ {
						assert.True(t, facts[i].UpdatedAt.After(facts[i+1].UpdatedAt) ||
							facts[i].UpdatedAt.Equal(facts[i+1].UpdatedAt),
							"Facts should be ordered by updated_at DESC")
					}
				}
			}
		})
	}
}
