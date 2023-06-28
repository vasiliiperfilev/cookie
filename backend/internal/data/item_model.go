package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
)

var ItemUnitsToId = map[string]int64{
	"l":  1,
	"kg": 2,
}

var IdToItemUnits = reverseMap(ItemUnitsToId)

type ItemModel interface {
	Insert(item *Item) error // TODO: use value instead of pointers
	GetById(id int64) (Item, error)
	GetAllBySupplierId(id int64) ([]Item, error)
	Update(item Item) (Item, error)
}

type PsqlItemModel struct {
	db *sql.DB
}

func NewPsqlItemModel(db *sql.DB) *PsqlItemModel {
	return &PsqlItemModel{db: db}
}

func (m PsqlItemModel) Insert(item *Item) error {
	query := `
    INSERT INTO items(supplier_id, unit_id, size, name, image_url)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING item_id
	`

	unitId, ok := ItemUnitsToId[item.Unit]
	if !ok {
		return ErrUnprocessableEntity
	}
	args := []any{item.SupplierId, unitId, item.Size, item.Name, item.ImageUrl}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, args...).Scan(&item.Id)
	if err != nil {
		return err
	}

	return nil
}

func (m PsqlItemModel) GetById(id int64) (Item, error) {
	query := `
		SELECT item_id, supplier_id, unit_id, size, name, image_url
		FROM items
		WHERE item_id=$1
	`

	var item Item
	var unitId int64

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, id).Scan(
		&item.Id,
		&item.SupplierId,
		&unitId,
		&item.Size,
		&item.Name,
		&item.ImageUrl,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Item{}, ErrRecordNotFound
		default:
			return Item{}, err
		}
	}
	item.Unit = IdToItemUnits[unitId]

	return item, nil
}

func (m PsqlItemModel) GetAllBySupplierId(id int64) ([]Item, error) {
	query := `
		SELECT item_id, supplier_id, unit_id, size, name, image_url
		FROM items
		WHERE supplier_id=$1
	`

	var items []Item = []Item{}

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
		item := Item{}
		var unitId int64
		if err := rows.Scan(
			&item.Id,
			&item.SupplierId,
			&unitId,
			&item.Size,
			&item.Name,
			&item.ImageUrl,
		); err != nil {
			return nil, err
		}
		item.Unit = IdToItemUnits[unitId]
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (m PsqlItemModel) Update(item Item) (Item, error) {
	query := `
		UPDATE items
		SET unit_id = $1, size = $2, name = $3, image_url = $4
		WHERE item_id = $5
	`
	unitId, ok := ItemUnitsToId[item.Unit]
	if !ok {
		return Item{}, ErrUnprocessableEntity
	}

	args := []any{unitId, item.Size, item.Name, item.ImageUrl, item.Id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return Item{}, err
	}

	return item, nil
}
