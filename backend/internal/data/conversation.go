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
	Id               int64             `json:"id"`
	UserIds          []int64           `json:"userIds"`
	LastMessageId    int64             `json:"lastMessageId"`
	LastReadMessages []LastReadMessage `json:"lastReadMessages"`
	Version          int               `json:"version"`
}

type LastReadMessage struct {
	UserId  int64   `json:"userId"`
	Message Message `json:"message"`
}

func AssertConversation(t *testing.T, got Conversation, want Conversation) {
	t.Helper()
	tester.AssertValue(t, got.Id, want.Id, "Expected same conversation id")
	tester.AssertValue(t, got.LastMessageId, want.LastMessageId, "Expected same last message id")
	if !EqualArrays(got.UserIds, want.UserIds) {
		t.Fatalf("Expected same userIds.Got %v, want %v", got.UserIds, want.UserIds)
	}
	if !EqualArrays(got.LastReadMessages, want.LastReadMessages) {
		t.Fatalf("Expected same userIds.Got %v, want %v", got.UserIds, want.UserIds)
	}
}
