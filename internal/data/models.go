package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	User IUserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		User: UserModel{DB: db},
	}
}
