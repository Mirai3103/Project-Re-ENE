package store

import (
	"github.com/jmoiron/sqlx"
)

type CharacterStore struct {
	db *sqlx.DB
}

func NewCharacterStore(db *sqlx.DB) *CharacterStore {
	return &CharacterStore{db: db}
}

func (c *CharacterStore) GetCharacter(id string) (*Character, error) {
	var character Character
	err := c.db.Get(&character, "SELECT * FROM characters WHERE id = ? LIMIT 1", id)
	if err != nil {
		return nil, cleanError(err)
	}
	return &character, nil
}

func (c *CharacterStore) GetCharacters() ([]*Character, error) {
	characters := []*Character{}
	err := c.db.Select(&characters, "SELECT * FROM characters")
	if err != nil {
		return nil, cleanError(err)
	}
	return characters, nil
}

func (c *CharacterStore) GetCharacterFacts(characterID string, limit int) ([]*CharacterFact, error) {
	characterFacts := []*CharacterFact{}
	err := c.db.Select(&characterFacts, "SELECT * FROM character_facts WHERE character_id = ? order by updated_at desc limit ?", characterID, limit)
	if err != nil {
		return nil, cleanError(err)
	}
	return characterFacts, nil
}
