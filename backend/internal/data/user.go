package data

import (
	"errors"
	"time"
	"unicode"

	"github.com/vasiliiperfilev/cookie/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	SupplierUserType  = 1
	ClientUserType    = 2
)

// add type enum
type User struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  password  `json:"-"`
	Type      int       `json:"type"`
	ImageId   string    `json:"imageId"`
	Version   int       `json:"-"`
}

type PostUserDto struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Type     int    `json:"type"`
	ImageId  string `json:"imageId"`
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

func ValidateRegisterUserInput(v *validator.Validator, input *PostUserDto) {
	ValidateEmail(v, input.Email)

	ValidatePasswordPlaintext(v, input.Password)

	v.Check(verifyUserType(input.Type), "type", "must be a valid user type")

	v.Check(input.Name != "", "name", "must be provided")

	v.Check(input.ImageId != "", "imageId", "must be provided")
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be an email")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(verifyNewPassword(password), "password", "must be at least 8 chars, have a special symbol, number, lower and upper case letter")
}

func verifyNewPassword(s string) bool {
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
	userTypes := map[string]int{"supplier": SupplierUserType, "client": ClientUserType}
	for _, value := range userTypes {
		if value == i {
			return true
		}
	}
	return false
}
