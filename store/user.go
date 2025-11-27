package store

import "github.com/jmoiron/sqlx"

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{db: db}
}

func (u *UserStore) GetUser(id string) (*User, error) {
	var user User
	err := u.db.Get(&user, "SELECT * FROM users WHERE id = ? LIMIT 1", id)
	if err != nil {
		return nil, cleanError(err)
	}
	return &user, nil
}

func (u *UserStore) GetUserFacts(userID string, limit int) ([]*UserFact, error) {
	userFacts := []*UserFact{}
	err := u.db.Select(&userFacts, "SELECT * FROM user_facts WHERE user_id = ? order by updated_at desc limit ?", userID, limit)
	if err != nil {
		return nil, cleanError(err)
	}
	return userFacts, nil
}
