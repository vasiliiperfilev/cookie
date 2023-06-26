package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
)

var ItemUnitsToId = map[string]int64{
	"l": 1,
}

var IdToItemUnits = reverseMap(ItemUnitsToId)

type ItemModel interface {
	Insert(item *Item) error
	GetById(id int64) (Item, error)
}

type PsqlItemModel struct {
	db *sql.DB
}

func NewPsqlItemModel(db *sql.DB) *PsqlItemModel {
	return &PsqlItemModel{db: db}
}

func (m PsqlItemModel) Insert(item *Item) error {
	query := `
        INSERT INTO items(supplier_id, unit_id, name, image_url)
        VALUES ($1, $2, $3, $4)
        RETURNING item_id`

	args := []any{item.SupplierId, ItemUnitsToId[item.Unit], item.Name, item.ImageUrl}

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
		SELECT item_id, supplier_id, unit_id, name, image_url
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
		unitId,
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
