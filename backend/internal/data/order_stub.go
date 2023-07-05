package data

type StubOrderModel struct {
	orders  map[int64]Order
	idCount int64
}

func NewStubOrderModel(orders []Order) *StubOrderModel {
	ordersMap := map[int64]Order{}
	for _, order := range orders {
		ordersMap[order.Id] = order
	}
	return &StubOrderModel{orders: ordersMap, idCount: int64(len(orders))}
}

func (s *StubOrderModel) Insert(order Order) (Order, error) {
	s.idCount++
	order.Id = s.idCount
	s.orders[order.Id] = order
	return order, nil
}

func (s *StubOrderModel) GetById(id int64) (Order, error) {
	if order, ok := s.orders[id]; !ok {
		return Order{}, ErrRecordNotFound
	} else {
		return order, nil
	}
}
