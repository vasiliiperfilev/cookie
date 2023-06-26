package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	User         UserModel
	Token        TokenModel
	Conversation ConversationModel
	Message      MessageModel
	Item         ItemModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		User:         NewPsqlUserModel(db),
		Token:        NewPsqlTokenModel(db),
		Conversation: NewPsqlConversationModel(db),
		Message:      NewPsqlMessageModel(db),
		Item:         NewPsqlItemModel(db),
	}
}
