package data

import (
	"context"
	"database/sql"
	"time"
)

type PermissionModel interface {
	GetAllForType(typeId int64) (Permissions, error)
}

type PsqlPermissionModel struct {
	db *sql.DB
}

func NewPsqlPermissionModel(db *sql.DB) *PsqlPermissionModel {
	return &PsqlPermissionModel{db: db}
}

func (m PsqlPermissionModel) GetAllForType(typeId int64) (Permissions, error) {
	query := `
        SELECT permission_id
        FROM types_permissions
        WHERE user_type_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, query, typeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions Permissions

	for rows.Next() {
		var permission int

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}
