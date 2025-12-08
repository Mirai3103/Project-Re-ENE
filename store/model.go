package store

import "time"

type Conversation struct {
	ID             string    `json:"id" db:"id"`
	Title          *string   `json:"title" db:"title"`
	MaxWindowSize  int       `json:"max_window_size" db:"max_window_size"`
	CharacterID    string    `json:"character_id" db:"character_id"`
	UserID         string    `json:"user_id" db:"user_id"`
	CurrentSummary *string   `json:"current_summary" db:"current_summary"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type ConversationMessage struct {
	ID             string    `json:"id" db:"id"`
	ConversationID string    `json:"conversation_id" db:"conversation_id"`
	Role           string    `json:"role" db:"role"`
	Content        string    `json:"content" db:"content"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type Character struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	BasePrompt  string    `json:"base_prompt" db:"base_prompt"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CharacterFact struct {
	ID          string    `json:"id" db:"id"`
	CharacterID string    `json:"character_id" db:"character_id"`
	Name        string    `json:"name" db:"name"`
	Value       string    `json:"value" db:"value"`
	Type        string    `json:"type" db:"type"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type User struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Bio       string    `json:"bio" db:"bio"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserFact struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Value     string    `json:"value" db:"value"`
	Type      string    `json:"type" db:"type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Memory struct {
	ID             string    `json:"id" db:"id"`
	UserID         string    `json:"user_id" db:"user_id"`
	CharacterID    string    `json:"character_id" db:"character_id"`
	Content        string    `json:"content" db:"content"`
	Importance     float64   `json:"importance" db:"importance"`
	LastAccessedAt time.Time `json:"last_accessed_at" db:"last_accessed_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
