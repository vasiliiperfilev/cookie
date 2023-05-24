package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"unicode"

	"github.com/vasiliiperfilev/cookie/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrEditConflict   = errors.New("edit conflict")
)

// add type enum
type User struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Type      int       `json:"type"`
	ImageId   string    `json:"imageId"`
	Version   int       `json:"-"`
}

type RegisterUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Type     int    `json:"type"`
	ImageId  string `json:"imageId"`
}

type RegisterUserResponse struct {
	Email   string `json:"email"`
	Type    int    `json:"type"`
	ImageId string `json:"imageId"`
}

// The plaintext field is a *pointer* to a string,
// so that we're able to distinguish between a plaintext password not being present in
// the struct at all, versus a plaintext password which is the empty string "".
type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateRegisterUserInput(v *validator.Validator, input *RegisterUserInput) {
	v.Check(input.Email != "", "email", "must be provided")
	v.Check(validator.Matches(input.Email, validator.EmailRX), "email", "must be an email")

	v.Check(verifyPassword(input.Password), "password", "must be at least 8 chars, have a special symbol, number, lower and upper case letter")

	v.Check(verifyUserType(input.Type), "type", "must be a valid user type")

	v.Check(input.ImageId != "", "imageId", "must be provided")
}

func verifyPassword(s string) bool {
	number := false
	upper := false
	lower := false
	special := false
	letters := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c):
			lower = true
		default:
			return false
		}
		letters++
	}
	// max bccrypt input size is 72 bytes
	size := letters >= 8 && letters <= 72
	return number && upper && lower && special && size
}

func verifyUserType(i int) bool {
	userTypes := map[string]int{"supplier": 1, "business": 2}
	for _, value := range userTypes {
		if value == i {
			return true
		}
	}
	return false
}

type IUserModel interface {
	Insert(user *User) error
	GetByEmail(email string) (*User, error)
	Update(user *User) error
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {
	query := `
        INSERT INTO users (name, email, password_hash) 
        VALUES ($1, $2, $3)
        RETURNING id, created_at, version`

	args := []any{user, user.Email, user.Password.hash}

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

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
        SELECT id, created_at, name, email, password_hash, activated, version
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
func (m UserModel) Update(user *User) error {
	query := `
        UPDATE users 
        SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
        WHERE id = $5 AND version = $6
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
