package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"
)

type UserModel interface {
	Insert(user *User) error
	GetByEmail(email string) (User, error)
	GetById(id int64) (User, error)
	Update(user User) error
	GetForToken(tokenScope, tokenPlaintext string) (User, error)
}

type PsqlUserModel struct {
	db *sql.DB
}

func NewPsqlUserModel(db *sql.DB) *PsqlUserModel {
	return &PsqlUserModel{db: db}
}

func (m PsqlUserModel) Insert(user *User) error {
	query := `
        INSERT INTO users (email, name, password_hash, user_type_id, image_id) 
        VALUES ($1, $2, $3, $4, $5)
        RETURNING user_id, created_at, version`

	args := []any{user.Email, user.Name, user.Password.hash, user.Type, user.ImageId}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// We check for a violation of the UNIQUE "users_email_key"
	// specifically, and return custom ErrDuplicateEmail error instead.
	err := m.db.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.CreatedAt, &user.Version)
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

func (m PsqlUserModel) GetByEmail(email string) (User, error) {
	query := `
        SELECT user_id, created_at, email, name, password_hash, user_type_id, version
        FROM users
        WHERE email = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.Email,
		&user.Name,
		&user.Password.hash,
		&user.Type,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return User{}, ErrRecordNotFound
		default:
			return User{}, err
		}
	}

	return user, nil
}

func (m PsqlUserModel) GetById(id int64) (User, error) {
	query := `
        SELECT user_id, created_at, email, name, password_hash, user_type_id, version
        FROM users
        WHERE user_id = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, id).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.Email,
		&user.Name,
		&user.Password.hash,
		&user.Type,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return User{}, ErrRecordNotFound
		default:
			return User{}, err
		}
	}

	return user, nil
}

// Update the details for a specific user. Notice that we check against the version
// field to help prevent any race conditions during the request cycle, just like we did
// when updating a movie. And we also check for a violation of the "users_email_key"
// constraint when performing the update, just like we did when inserting the user
// record originally.
func (m PsqlUserModel) Update(user User) error {
	query := `
        UPDATE users
        SET email = $1, password_hash = $2, name = $3, version = version + 1
        WHERE user_id = $4 AND version = $5
        RETURNING version`

	args := []any{
		user.Email,
		user.Password.hash,
		user.Name,
		user.Id,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, args...).Scan(&user.Version)
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

func (m PsqlUserModel) GetForToken(tokenScope, tokenPlaintext string) (User, error) {
	// Calculate the SHA-256 hash of the plaintext token provided by the client.
	// Remember that this returns a byte *array* with length 32, not a slice.
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
        SELECT a.user_id, a.created_at, a.email, a.name, a.password_hash, a.user_type_id, a.image_id, a.version
        FROM users as a
        INNER JOIN tokens as t
        ON a.user_id = t.user_id
        WHERE t.hash = $1
        AND t.scope = $2 
        AND t.expiry > $3`

	// Create a slice containing the query arguments. Notice how we use the [:] operator
	// to get a slice containing the token hash, rather than passing in the array (which
	// is not supported by the pq driver)
	args := []any{tokenHash[:], tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query, scanning the return values into a User struct. If no matching
	// record is found we return an ErrRecordNotFound error.
	err := m.db.QueryRowContext(ctx, query, args...).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.Email,
		&user.Name,
		&user.Password.hash,
		&user.Type,
		&user.ImageId,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return User{}, ErrRecordNotFound
		default:
			return User{}, err
		}
	}

	// Return the matching user.
	return user, nil
}
