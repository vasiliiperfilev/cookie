package data

import (
	"context"
	"database/sql"
	"time"
)

type OrderRepository interface {
	Insert(dto PostOrderDto) (Order, error)
}

type PsqlOrderRepository struct {
	db      *sql.DB
	order   OrderModel
	message MessageModel
}

func NewPsqlOrderRepository(db *sql.DB, order OrderModel, message MessageModel) PsqlOrderRepository {
	return PsqlOrderRepository{order: order, message: message, db: db}
}

func (r PsqlOrderRepository) Insert(dto PostOrderDto) (Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Order{}, err
	}
	defer tx.Rollback()

	// check if user in conversation

	// Create a message to store order
	message := Message{
		ConversationId: dto.ConversationId,
		Content:        "Order created",
		SenderId:       dto.ClientId,
	}
	err = r.message.Insert(&message)
	if err != nil {
		return Order{}, err
	}

	// Create an order linked to the message
	order := Order{
		Items:     dto.Items,
		StateId:   OrderStateCreated,
		MessageId: message.Id,
	}
	order, err = r.order.Insert(order)
	if err != nil {
		return Order{}, err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return Order{}, err
	}

	return order, nil
}
