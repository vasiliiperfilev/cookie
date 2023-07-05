package data

type OrderModel interface {
	Insert(order Order) (Order, error)
	GetById(id int64) (Order, error)
}
