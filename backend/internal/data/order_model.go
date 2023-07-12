package data

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/lib/pq"
)

type OrderModel interface {
	Insert(order Order) (Order, error)
	GetById(id int64) (Order, error)
	GetAllByUserId(id int64) ([]Order, error)
	Update(order Order) (Order, error)
}

type PsqlOrderModel struct {
	db *sql.DB
}

func NewPsqlOrderModel(db *sql.DB) *PsqlOrderModel {
	return &PsqlOrderModel{db: db}
}

func (m PsqlOrderModel) Insert(order Order) (Order, error) {
	query := `
    INSERT INTO orders(message_id)
    VALUES ($1)
    RETURNING order_id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, order.MessageId).Scan(&order.Id, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return Order{}, err
	}

	query = `
    INSERT INTO orders_states(order_id, state_id)
    VALUES ($1, $2)
	`
	args := []any{order.Id, order.StateId}

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return Order{}, err
	}
	defer rows.Close()

	txn, err := m.db.Begin()
	if err != nil {
		return Order{}, err
	}

	stmt, err := txn.Prepare(pq.CopyIn("orders_items", "order_id", "item_id", "quantity"))
	if err != nil {
		return Order{}, err
	}

	for _, iq := range order.Items {
		_, err = stmt.Exec(order.Id, iq.ItemId, iq.Quantity)
		if err != nil {
			return Order{}, err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return Order{}, err
	}

	err = stmt.Close()
	if err != nil {
		return Order{}, err
	}

	err = txn.Commit()
	if err != nil {
		return Order{}, err
	}

	return order, nil
}

func (m PsqlOrderModel) GetById(id int64) (Order, error) {
	if id < 1 {
		return Order{}, ErrRecordNotFound
	}
	query := `
		SELECT o.order_id, o.message_id, o.created_at, o.updated_at, os.state_id, json_agg(json_build_object(
			'item_id', oi.item_id, 
			'quantity', oi.quantity
		)) as items
		FROM orders as o
			INNER JOIN orders_items as oi ON o.order_id = oi.order_id
			INNER JOIN orders_states as os ON o.order_id = os.order_id
		WHERE o.order_id=$1
			AND os.created_at = (
				SELECT MAX(created_at) FROM orders_states WHERE order_id=$1
			)
		GROUP BY o.order_id, os.state_id
	`

	var order Order

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var items []byte
	err := m.db.QueryRowContext(ctx, query, id).Scan(
		&order.Id,
		&order.MessageId,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.StateId,
		&items,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Order{}, ErrRecordNotFound
		default:
			return Order{}, err
		}
	}

	err = json.NewDecoder(bytes.NewBuffer(items)).Decode(&order.Items)
	if err != nil {
		return Order{}, err
	}

	return order, nil
}

func (m PsqlOrderModel) GetAllByUserId(id int64) ([]Order, error) {
	if id < 1 {
		return []Order{}, ErrRecordNotFound
	}
	query := `
		SELECT o.order_id, o.message_id, o.created_at, o.updated_at, os.state_id, json_agg(json_build_object(
			'item_id', oi.item_id, 
			'quantity', oi.quantity
		)) as items
		FROM orders as o
			INNER JOIN orders_items as oi ON o.order_id = oi.order_id
			INNER JOIN messages as m ON m.message_id = o.message_id
			INNER JOIN conversations as c ON c.conversation_id = m.conversation_id
			INNER JOIN conversations_users as cu ON cu.conversation_id = c.conversation_id
			INNER JOIN (
				SELECT 
					order_id, 
					state_id, 
					rank() over (partition by order_id order by created_at desc) as rownum
				FROM orders_states
			) as os ON o.order_id = os.order_id
		WHERE cu.user_id=$1 AND os.rownum = 1
		GROUP BY o.order_id, os.state_id
	`

	orders := []Order{}

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
		order := Order{}
		var items []byte
		if err := rows.Scan(
			&order.Id,
			&order.MessageId,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.StateId,
			&items,
		); err != nil {
			return nil, err
		}
		err = json.NewDecoder(bytes.NewBuffer(items)).Decode(&order.Items)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (m PsqlOrderModel) Update(order Order) (Order, error) {
	if order.Id < 1 {
		return Order{}, ErrRecordNotFound
	}

	prevOrder, err := m.GetById(order.Id)
	if err != nil {
		return Order{}, err
	}

	txn, err := m.db.Begin()
	if err != nil {
		return Order{}, err
	}

	query := `
		UPDATE orders
		SET message_id = $1
		WHERE order_id = $2
	`

	args := []any{order.MessageId, order.Id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := txn.QueryContext(ctx, query, args...)
	if err != nil {
		return Order{}, err
	}
	rows.Close()
	// update state if required
	if order.StateId != prevOrder.StateId {
		query = `
    	INSERT INTO orders_states(order_id, state_id)
    	VALUES ($1, $2)
		`
		args := []any{order.Id, order.StateId}

		rows, err := txn.QueryContext(ctx, query, args...)
		if err != nil {
			return Order{}, err
		}
		rows.Close()
	}

	// update items if required
	if !EqualArrays(order.Items, prevOrder.Items) {

		for _, iq := range prevOrder.Items {
			query := `
        DELETE FROM orders_items
        WHERE order_id = $1 AND item_id = $2`

			_, err := txn.Exec(query, prevOrder.Id, iq.ItemId)
			if err != nil {
				return Order{}, err
			}
		}

		stmt, err := txn.Prepare(pq.CopyIn("orders_items", "order_id", "item_id", "quantity"))
		if err != nil {
			return Order{}, err
		}

		for _, iq := range order.Items {
			_, err = stmt.Exec(order.Id, iq.ItemId, iq.Quantity)
			if err != nil {
				return Order{}, err
			}
		}

		_, err = stmt.Exec()
		if err != nil {
			return Order{}, err
		}

		err = stmt.Close()
		if err != nil {
			return Order{}, err
		}
	}

	err = txn.Commit()
	if err != nil {
		return Order{}, err
	}

	return order, nil
}
