package data

type Conversation struct {
	Id            int64   `json:"id"`
	UserIds       []int64 `json:"userIds"`
	LastMessageId int64   `json:"lastMessageId"`
}
