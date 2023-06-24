package data

type ItemModel interface {
	Insert(item *Item) error
	GetById(id int64) (Item, error)
}
