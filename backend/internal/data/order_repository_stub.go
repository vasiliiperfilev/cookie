package data

type StubOrderRepository struct {
	order   *StubOrderModel
	message *StubMessageModel
}

func NewStubOrderRepository(order *StubOrderModel, message *StubMessageModel) *StubOrderRepository {
	return &StubOrderRepository{order: order, message: message}
}

func (r *StubOrderRepository) Insert(dto PostOrderDto) (Order, error) {
	// Create a message to store order
	message := Message{
		ConversationId: dto.ConversationId,
		Content:        "Order created",
		SenderId:       dto.ClientId,
	}
	err := r.message.Insert(&message)
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
		// delete created message to imitate transaction
		r.message.DeleteById(message.Id)
		return Order{}, err
	}

	return order, nil
}
