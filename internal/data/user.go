package data

import (
	"unicode"

	"github.com/vasiliiperfilev/cookie/internal/validator"
)

// add type enum
type User struct {
	Id                int    `json:"id"`
	Email             string `json:"email"`
	PasswordHash      string `json:"-"`
	Type              string `json:"type"`
	ImageId           string `json:"imageId"`
	LastReadMessageId int    `json:"-"`
}

type RegisterUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Type     string `json:"type"`
	ImageId  string `json:"imageId"`
}

type RegisterUserResponse struct {
	Email   string `json:"email"`
	Type    string `json:"type"`
	ImageId string `json:"imageId"`
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
	size := letters >= 8
	return number && upper && lower && special && size
}

func verifyUserType(s string) bool {
	userTypes := []string{"supplier", "business"}
	for _, value := range userTypes {
		if value == s {
			return true
		}
	}
	return false
}
