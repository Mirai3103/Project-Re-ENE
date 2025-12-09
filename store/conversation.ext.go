package store

import (
	"context"
	"database/sql"
	"errors"
)

//sql.ErrNoRows

func (q Queries) CreateConversationIfNotExists(ctx context.Context, arg CreateConversationParams) (Conversation, error) {
	exist, err := q.GetConversation(ctx, arg.ID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		err = q.CreateConversation(ctx, arg)
		if err != nil {
			return Conversation{}, err
		}
		exist, err = q.GetConversation(ctx, arg.ID)
	}

	return exist, err
}
