package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type MessageModel interface {
	Insert(msg *Message) error // TODO: use value, not pointer
	GetAllByConversationId(id int64) ([]Message, error)
	GetById(id int64) (Message, error)
	Update(msg Message) error
}

type PsqlMessageModel struct {
	db *sql.DB
}

func NewPsqlMessageModel(db *sql.DB) *PsqlMessageModel {
	return &PsqlMessageModel{db: db}
}

func (m PsqlMessageModel) Insert(msg *Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = insertMessage(msg, tx)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func insertMessage(msg *Message, tx *sql.Tx) error {
	query := `
        INSERT INTO messages(sender_id, prev_message_id, conversation_id, content)
        VALUES ($1, $2, $3, $4)
        RETURNING message_id, created_at`

	args := []any{msg.SenderId, msg.PrevMessageId, msg.ConversationId, msg.Content}

	err := tx.QueryRow(query, args...).Scan(&msg.Id, &msg.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (m PsqlMessageModel) GetAllByConversationId(id int64) ([]Message, error) {
	query := `
	    SELECT message_id, sender_id, conversation_id, prev_message_id, created_at, content
	    FROM messages
	    WHERE conversation_id = $1`

	messages := make([]Message, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, query, id)
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
		msg := Message{}
		if err := rows.Scan(&msg.Id, &msg.SenderId, &msg.ConversationId, &msg.PrevMessageId, &msg.CreatedAt, &msg.Content); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (m PsqlMessageModel) GetById(id int64) (Message, error) {
	query := `
	    SELECT message_id, sender_id, conversation_id, prev_message_id, created_at, content
	    FROM messages
	    WHERE message_id = $1`

	msg := Message{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, id).Scan(&msg.Id, &msg.SenderId, &msg.ConversationId, &msg.PrevMessageId, &msg.CreatedAt, &msg.Content)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Message{}, ErrRecordNotFound
		default:
			return Message{}, err
		}
	}

	return msg, nil
}

func (m PsqlMessageModel) Update(msg Message) error {
	if msg.Id < 1 {
		return ErrRecordNotFound
	}
	query := `
			UPDATE messages
			SET content = $1
			WHERE message_id = $2
		`

	args := []any{msg.Content, msg.Id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
