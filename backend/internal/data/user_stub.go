package data

import (
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type StubUserModel struct {
	users   []User
	idCount int64
}

func NewStubUserModel(users []User) *StubUserModel {
	return &StubUserModel{users: users}
}

func (s *StubUserModel) Insert(user *User) error {
	if _, err := s.GetByEmail(user.Email); err == nil {
		return ErrDuplicateEmail
	}
	s.idCount++
	user.Id = s.idCount
	s.users = append(s.users, *user)
	return nil
}

func (s *StubUserModel) GetByEmail(email string) (User, error) {
	for _, user := range s.users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, ErrRecordNotFound
}

func (s *StubUserModel) GetById(id int64) (User, error) {
	for _, user := range s.users {
		if user.Id == id {
			return user, nil
		}
	}
	return User{}, ErrRecordNotFound
}

func (s *StubUserModel) Update(user User) error {
	if _, err := s.GetByEmail(user.Email); err == nil {
		return ErrDuplicateEmail
	}
	for k, v := range s.users {
		if v.Id == user.Id {
			s.users[k] = user
		}
	}
	return nil
}

// takes first symbol from plaintext token and use it as Id to find a user
func (s *StubUserModel) GetForToken(tokenScope, tokenPlaintext string) (User, error) {
	id, err := strconv.ParseInt(string(tokenPlaintext[0]), 10, 64)
	if err != nil {
		return User{}, err
	}
	for _, user := range s.users {
		if user.Id == id {
			return user, nil
		}
	}
	return User{}, ErrRecordNotFound
}

func (s *StubUserModel) GetAllBySearch(query string) ([]User, error) {
	result := []User{}
	queryNoSymbols := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(query, "")
	keywords := strings.Split(queryNoSymbols, " ")
	for _, user := range s.users {
		nameNoSymbols := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(user.Name, "")
		nameWords := strings.Split(nameNoSymbols, " ")
		for _, keyword := range keywords {
			if slices.Contains(nameWords, keyword) {
				result = append(result, user)
			}
		}
	}
	return result, nil
}
