package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type ConversationModel interface {
	Insert(conversation PostConversationDto) (Conversation, error)
	GetAllByUserId(userId int64) ([]Conversation, error)
	GetById(id int64) (Conversation, error)
}

type PostConversationDto struct {
	UserIds []int64 `json:"userIds"`
}

type PsqlConversationModel struct {
	db *sql.DB
}

func NewPsqlConversationModel(db *sql.DB) *PsqlConversationModel {
	return &PsqlConversationModel{db: db}
}

func (m PsqlConversationModel) Insert(dto PostConversationDto) (Conversation, error) {
	query := `
        INSERT INTO conversations(last_message_id)
        VALUES (0)
        RETURNING conversation_id, version`

	cvs := Conversation{
		UserIds: dto.UserIds,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query).Scan(&cvs.Id, &cvs.Version)
	if err != nil {
		return Conversation{}, err
	}

	err = m.insertConversationUsers(&cvs)
	if err != nil {
		// remove empty conversation
		return Conversation{}, err
	}

	return cvs, nil
}

func (m PsqlConversationModel) GetAllByUserId(userId int64) ([]Conversation, error) {
	query := `
    SELECT c.conversation_id, c.last_message_id, c.version, array_agg(c_u_ids.user_id) as user_ids
    FROM conversations_users as c_u
			INNER JOIN conversations as c
				ON c.conversation_id = c_u.conversation_id
			INNER JOIN conversations_users as c_u_ids
				ON c.conversation_id = c_u_ids.conversation_id
    WHERE c_u.user_id = $1
		GROUP BY c.conversation_id`

	conversations := []Conversation{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, query, userId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		conversation := Conversation{}
		if err := rows.Scan(&conversation.Id, &conversation.LastMessageId, &conversation.Version, (*pq.Int64Array)(&conversation.UserIds)); err != nil {
			return nil, err
		}
		conversations = append(conversations, conversation)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}

func (m PsqlConversationModel) GetById(id int64) (Conversation, error) {
	query := `
		SELECT c.conversation_id, c.last_message_id, c.version, array_agg(c_u.user_id) as user_ids
		FROM conversations_users as c_u
			INNER JOIN conversations as c
				ON c.conversation_id = c_u.conversation_id
		WHERE c.conversation_id = $1
		GROUP BY c.conversation_id`

	var conversation Conversation

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, id).Scan(
		&conversation.Id,
		&conversation.LastMessageId,
		&conversation.Version,
		(*pq.Int64Array)(&conversation.UserIds),
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Conversation{}, ErrRecordNotFound
		default:
			return Conversation{}, err
		}
	}

	return conversation, nil
}

func (m PsqlConversationModel) insertConversationUsers(conversation *Conversation) error {
	for _, userId := range conversation.UserIds {
		query := `
		INSERT INTO conversations_users(conversation_id, user_id)
		VALUES ($1, $2)`

		args := []any{conversation.Id, userId}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_, err := m.db.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}

	}
	return nil
}
