package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrEditConflict        = errors.New("edit conflict")
	ErrUnprocessableEntity = errors.New("can't process value")
)

type Models struct {
	User         UserModel
	Token        TokenModel
	Conversation ConversationModel
	Message      MessageModel
	Item         ItemModel
	Permission   PermissionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		User:         NewPsqlUserModel(db),
		Token:        NewPsqlTokenModel(db),
		Conversation: NewPsqlConversationModel(db),
		Message:      NewPsqlMessageModel(db),
		Item:         NewPsqlItemModel(db),
		Permission:   NewPsqlPermissionModel(db),
	}
}
