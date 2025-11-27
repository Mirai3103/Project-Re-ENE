package store

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

var ErrNoRecordFound = errors.New("no record found")

type Store struct {
	*ConversationStore
	*CharacterStore
	*UserStore
}

func New(db *sqlx.DB) *Store {
	return &Store{
		NewConversationStore(db),
		NewCharacterStore(db),
		NewUserStore(db),
	}
}

func cleanError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNoRecordFound
	}
	return err
}
