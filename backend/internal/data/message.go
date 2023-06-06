package data

import "time"

type Message struct {
	Id             int64
	ConversationId int64
	Content        string
	SenderId       int64
	PrevMessageId  int64
	CreatedAt      time.Time
}
