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
	ClientId       int64
	ItemIds        []int64 `json:"itemIds"`
	ConversationId int64
}

type PatchOrderDto struct {
	ItemIds []int64 `json:"itemIds,omitempty"`
	StateId int     `json:"stateId,omitempty"`
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

func ValidatePostOrderInput(v *validator.Validator, dto PostOrderDto) {
	v.Check(len(dto.ItemIds) > 0, "itemIds", "must have at least 1 item")
}

func ValidatePatchOrderInput(v *validator.Validator, dto PatchOrderDto) {
	hasItems := len(dto.ItemIds) > 0
	validState := dto.StateId > 0 && dto.StateId <= 7
	if hasItems && validState {
		v.AddError("itemIds", "can't change both items and state")
		v.AddError("stateId", "can't change both items and state")
	}
	if !hasItems && !validState {
		v.AddError("itemIds", "valid items or state change is required")
		v.AddError("stateId", "valid items or state change is required")
	}
}
