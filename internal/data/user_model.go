package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type UserModel interface {
	Insert(user *User) error
	GetByEmail(email string) (*User, error)
	Update(user *User) error
}

type PsqlUserModel struct {
	DB *sql.DB
}

func (m PsqlUserModel) Insert(user *User) error {
	query := `
        INSERT INTO app_user (email, password_hash, user_type_id) 
        VALUES ($1, $2, $3)
        RETURNING user_id, created_at, version`

	args := []any{user.Email, user.Password.hash, user.Type}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// We check for a violation of the UNIQUE "users_email_key"
	// specifically, and return custom ErrDuplicateEmail error instead.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m PsqlUserModel) GetByEmail(email string) (*User, error) {
	query := `
        SELECT user_id, created_at, email, password_hash, user_type_id, version
        FROM users
        WHERE email = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.Email,
		&user.Password.hash,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

// Update the details for a specific user. Notice that we check against the version
// field to help prevent any race conditions during the request cycle, just like we did
// when updating a movie. And we also check for a violation of the "users_email_key"
// constraint when performing the update, just like we did when inserting the user
// record originally.
func (m PsqlUserModel) Update(user *User) error {
	query := `
        UPDATE users 
        SET email = $1, password_hash = $2, version = version + 1
        WHERE id = $3 AND version = $4
        RETURNING version`

	args := []any{
		user.Email,
		user.Password.hash,
		user.Id,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}
