package data

import (
	"errors"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/tester"
)

var (
	ErrDuplicateConversation = errors.New("duplicate conversation")
)

type Conversation struct {
	Id          int64   `json:"id"`
	Users       []User  `json:"users"`
	LastMessage Message `json:"lastMessage"`
	Version     int     `json:"version"`
}

func AssertConversation(t *testing.T, got Conversation, want Conversation) {
	t.Helper()
	tester.AssertValue(t, got.Id, want.Id, "Expected same conversation id")
	tester.AssertValue(t, got.LastMessage.Id, want.LastMessage.Id, "Expected same last message")
	getUserId := func(u User) int64 { return u.Id }
	gotIds := Map(got.Users, getUserId)
	wantIds := Map(got.Users, getUserId)
	if !EqualArraysContent(gotIds, wantIds) {
		t.Fatal("Expected same users")
	}
}
