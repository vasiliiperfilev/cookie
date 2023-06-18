package data

import "time"

type Message struct {
	Id             int64     `json:"id"`
	ConversationId int64     `json:"conversationId"`
	Content        string    `json:"content"`
	SenderId       int64     `json:"senderId"`
	PrevMessageId  int64     `json:"prevMessageId"`
	CreatedAt      time.Time `json:"createdAt"`
}

type PostMessageDto struct {
	ConversationId int64  `json:"conversationId"`
	Content        string `json:"content"`
	PrevMessageId  int64  `json:"prevMessageId"`
}
