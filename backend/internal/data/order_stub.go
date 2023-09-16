package data

import (
	"golang.org/x/exp/slices"
)

type StubOrderModel struct {
	orders       map[int64]Order
	conversation *StubConversationModel
	message      *StubMessageModel
	item         *StubItemModel
	idCount      int64
}

func NewStubOrderModel(orders []Order, item *StubItemModel, conversation *StubConversationModel, message *StubMessageModel) *StubOrderModel {
	ordersMap := map[int64]Order{}
	for _, order := range orders {
		ordersMap[order.Id] = order
	}
	return &StubOrderModel{
		orders:       ordersMap,
		idCount:      int64(len(orders)),
		conversation: conversation,
		message:      message,
		item:         item,
	}
}

func (s *StubOrderModel) Insert(dto PostOrderDto) (Order, error) {
	for _, item := range dto.Items {
		_, err := s.item.GetById(item.ItemId)
		if err != nil {
			return Order{}, ErrUnprocessableEntity
		}
	}
	msg := Message{
		ConversationId: dto.ConversationId,
		SenderId:       dto.ClientId,
		Content:        "Order created",
	}
	s.message.Insert(&msg)

	s.idCount++
	order := Order{
		MessageId: msg.Id,
		Items:     dto.Items,
		StateId:   OrderStateCreated,
	}
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

func (s *StubOrderModel) GetAllByUserId(id int64) ([]Order, error) {
	result := []Order{}
	for _, order := range s.orders {
		message, err := s.message.GetById(order.MessageId)
		if err != nil {
			return nil, ErrRecordNotFound
		}
		conversation, _ := s.conversation.GetById(message.ConversationId)
		if slices.ContainsFunc(conversation.Users, func(u User) bool { return u.Id == id }) {
			result = append(result, order)
		}
	}
	return result, nil
}

func (s *StubOrderModel) Update(order Order) (Order, error) {
	for _, item := range order.Items {
		_, err := s.item.GetById(item.ItemId)
		if err != nil {
			return Order{}, ErrUnprocessableEntity
		}
	}
	if _, ok := s.orders[order.Id]; !ok {
		return Order{}, ErrRecordNotFound
	} else {
		s.orders[order.Id] = order
	}
	return order, nil
}
