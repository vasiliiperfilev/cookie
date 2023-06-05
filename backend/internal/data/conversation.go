package data

import "errors"

var (
	ErrDuplicateConversation = errors.New("duplicate conversation")
)

type Conversation struct {
	Id            int64   `json:"id"`
	UserIds       []int64 `json:"userIds"`
	LastMessageId int64   `json:"lastMessageId"`
}
