package data

import (
	"time"
)

type StubTokenModel struct {
	tokens []Token
}

func NewStubTokenModel(tokens []Token) *StubTokenModel {
	return &StubTokenModel{tokens: tokens}
}

func (s *StubTokenModel) New(userID int64, ttl time.Duration, scope string) (Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return Token{}, err
	}

	s.tokens = append(s.tokens, token)
	return token, err
}

func (s *StubTokenModel) DeleteAllForUser(scope string, userID int64) error {
	for i, token := range s.tokens {
		if token.Scope == scope && token.UserId == userID {
			s.tokens = remove(s.tokens, i)
		}
	}
	return nil
}

func remove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
