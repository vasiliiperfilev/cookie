package data

type MessageModel interface {
	GetAllByUserId(id int64) ([]Message, error)
	Insert(msg Message) error
}
