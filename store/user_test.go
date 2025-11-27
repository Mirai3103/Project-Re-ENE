package store

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupUserTestData(t *testing.T, db *sqlx.DB) (string, string) {
	t.Helper()

	userID := uuid.New().String()
	userID2 := uuid.New().String()

	// Insert test users
	_, err := db.Exec(`
		INSERT INTO users (id, name, bio, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, userID, "Test User", "Test bio", time.Now(), time.Now())
	require.NoError(t, err)

	_, err = db.Exec(`
		INSERT INTO users (id, name, bio, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, userID2, "Another User", "Another bio", time.Now(), time.Now())
	require.NoError(t, err)

	// Insert test user facts
	for i := 0; i < 5; i++ {
		factID := uuid.New().String()
		_, err := db.Exec(`
			INSERT INTO user_facts (id, user_id, name, value, type, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, factID, userID, "fact"+string(rune(i)), "value"+string(rune(i)), "type"+string(rune(i)),
			time.Now().Add(time.Duration(i)*time.Second), time.Now().Add(time.Duration(i)*time.Second))
		require.NoError(t, err)
	}

	return userID, userID2
}

func TestUserStore_GetUser(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	store := NewUserStore(db)
	userID, _ := setupUserTestData(t, db)

	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{
			name:    "existing user",
			id:      userID,
			wantErr: nil,
		},
		{
			name:    "non-existing user",
			id:      "non-existing-id",
			wantErr: ErrNoRecordFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := store.GetUser(tt.id)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.id, user.ID)
				assert.Equal(t, "Test User", user.Name)
				assert.Equal(t, "Test bio", user.Bio)
			}
		})
	}
}

func TestUserStore_GetUserFacts(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	store := NewUserStore(db)
	userID, _ := setupUserTestData(t, db)

	tests := []struct {
		name          string
		userID        string
		limit         int
		wantMinCount  int
		wantMaxCount  int
		wantErr       error
		checkOrdering bool
	}{
		{
			name:          "get all facts",
			userID:        userID,
			limit:         10,
			wantMinCount:  5,
			wantMaxCount:  5,
			wantErr:       nil,
			checkOrdering: true,
		},
		{
			name:         "limit to 3 facts",
			userID:       userID,
			limit:        3,
			wantMinCount: 3,
			wantMaxCount: 3,
			wantErr:      nil,
		},
		{
			name:         "limit to 1 fact",
			userID:       userID,
			limit:        1,
			wantMinCount: 1,
			wantMaxCount: 1,
			wantErr:      nil,
		},
		{
			name:         "non-existing user",
			userID:       "non-existing-id",
			limit:        10,
			wantMinCount: 0,
			wantMaxCount: 0,
			wantErr:      nil, // Empty result is not an error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			facts, err := store.GetUserFacts(tt.userID, tt.limit)
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
