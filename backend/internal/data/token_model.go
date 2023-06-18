package data

import (
	"context"
	"database/sql"
	"time"
)

type TokenModel interface {
	New(userID int64, ttl time.Duration, scope string) (Token, error)
	DeleteAllForUser(scope string, userID int64) error
}

type PsqlTokenModel struct {
	db *sql.DB
}

func NewPsqlTokenModel(db *sql.DB) *PsqlTokenModel {
	return &PsqlTokenModel{db: db}
}

func (m PsqlTokenModel) New(userID int64, ttl time.Duration, scope string) (Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return Token{}, err
	}

	err = m.insert(token)
	return token, err
}

func (m PsqlTokenModel) insert(token Token) error {
	query := `
        INSERT INTO tokens (hash, user_id, expiry, scope) 
        VALUES ($1, $2, $3, $4)`

	args := []any{token.Hash, token.UserId, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

func (m PsqlTokenModel) DeleteAllForUser(scope string, userID int64) error {
	query := `
        DELETE FROM tokens
        WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, scope, userID)
	return err
}
