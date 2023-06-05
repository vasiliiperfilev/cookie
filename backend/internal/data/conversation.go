package data

import (
	"errors"
	"sort"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/tester"
)

var (
	ErrDuplicateConversation = errors.New("duplicate conversation")
)

type Conversation struct {
	Id            int64   `json:"id"`
	UserIds       []int64 `json:"userIds"`
	LastMessageId int64   `json:"lastMessageId"`
	Version       int     `json:"version"`
}

func AssertConversation(t *testing.T, got Conversation, want Conversation) {
	t.Helper()
	tester.AssertValue(t, want.Id, got.Id, "Expected same conversation id")
	tester.AssertValue(t, want.Id, got.Id, "Expected same last message id")
	sort.Slice(got.UserIds, func(i, j int) bool {
		return got.UserIds[i] >= got.UserIds[j]
	})
	sort.Slice(want.UserIds, func(i, j int) bool {
		return want.UserIds[i] >= want.UserIds[j]
	})
	tester.AssertValue(t, want.UserIds, got.UserIds, "Expected same usersId")
}
