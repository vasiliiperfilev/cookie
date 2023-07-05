package data

import (
	"time"

	"github.com/vasiliiperfilev/cookie/internal/validator"
)

type Order struct {
	Id        int64     `json:"id"`
	MessageId int64     `json:"messageId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ItemIds   []int64   `json:"itemIds"`
	StateId   int       `json:"stateId"`
}

type PostOrderDto struct {
	ItemIds        []int64 `json:"itemIds"`
	ConversationId int64
}

const (
	OrderStateCreated              = 1
	OrderStateAccepted             = 2
	OrderStateDeclined             = 3
	OrderStateFulfilled            = 4
	OrderStateConfirmedFulfillment = 5
	OrderStateSupplierChanges      = 6
	OrderStateClientChanges        = 7
)

func ValidatePostOrderInput(v *validator.Validator, order PostOrderDto) {
	v.Check(len(order.ItemIds) > 0, "itemIds", "must have at least 1 item")
}
