package data

import (
	"context"
	"database/sql"
	"time"
)

type TokenModel interface {
	New(userID int64, ttl time.Duration, scope string) (*Token, error)
	DeleteAllForUser(scope string, userID int64) error
}

type PsqlTokenModel struct {
	DB *sql.DB
}

func (m PsqlTokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.insert(token)
	return token, err
}

func (m PsqlTokenModel) insert(token *Token) error {
	query := `
        INSERT INTO token (hash, app_user_id, expiry, scope) 
        VALUES ($1, $2, $3, $4)`

	args := []any{token.Hash, token.UserId, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m PsqlTokenModel) DeleteAllForUser(scope string, userID int64) error {
	query := `
        DELETE FROM token 
        WHERE scope = $1 AND app_user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, scope, userID)
	return err
}
