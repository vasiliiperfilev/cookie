package data

type ConversationModel interface {
	Insert(conversation *Conversation) error
	// rename to get all
	GetAllByUserId(userId int64) ([]*Conversation, error)
}
