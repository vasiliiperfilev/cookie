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

	cvs := Conversation{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query).Scan(&cvs.Id, &cvs.Version)
	if err != nil {
		return Conversation{}, err
	}

	err = m.insertConversationUsers(&cvs, dto.UserIds)
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
		userIds := []int64{}
		var lastMessageId int64
		if err := rows.Scan(&conversation.Id, &lastMessageId, &conversation.Version, (*pq.Int64Array)(&userIds)); err != nil {
			return nil, err
		}
		err = m.getUsers(&conversation, userIds)
		if err != nil {
			return nil, err
		}
		err = m.getLastMessage(&conversation, lastMessageId)
		if err != nil {
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
	var lastMessageId int64
	userIds := []int64{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, id).Scan(
		&conversation.Id,
		&lastMessageId,
		&conversation.Version,
		(*pq.Int64Array)(&userIds),
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Conversation{}, ErrRecordNotFound
		default:
			return Conversation{}, err
		}
	}

	err = m.getUsers(&conversation, userIds)
	if err != nil {
		return Conversation{}, err
	}
	err = m.getLastMessage(&conversation, lastMessageId)
	if err != nil {
		return Conversation{}, err
	}

	return conversation, nil
}

func (m PsqlConversationModel) insertConversationUsers(conversation *Conversation, userIds []int64) error {
	for _, userId := range userIds {
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

func (m PsqlConversationModel) getLastMessage(conversation *Conversation, messageId int64) error {
	query := `
	    SELECT message_id, sender_id, conversation_id, prev_message_id, created_at, content
	    FROM messages
	    WHERE message_id = $1`

	msg := Message{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, messageId).Scan(&msg.Id, &msg.SenderId, &msg.ConversationId, &msg.PrevMessageId, &msg.CreatedAt, &msg.Content)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	conversation.LastMessage = msg
	return nil
}

func (m PsqlConversationModel) getUsers(conversation *Conversation, userIds []int64) error {
	users := []User{}
	for _, id := range userIds {
		query := `
        SELECT user_id, created_at, email, name, password_hash, user_type_id, version, image_id
        FROM users
        WHERE user_id = $1`

		var user User

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := m.db.QueryRowContext(ctx, query, id).Scan(
			&user.Id,
			&user.CreatedAt,
			&user.Email,
			&user.Name,
			&user.Password.hash,
			&user.Type,
			&user.Version,
			&user.ImageId,
		)

		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return ErrRecordNotFound
			default:
				return err
			}
		}

		users = append(users, user)
	}

	conversation.Users = users
	return nil
}
