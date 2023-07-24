package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/vasiliiperfilev/cookie/internal/validator"
)

type ItemQuantity struct {
	ItemId   int64 `json:"itemId"`
	Quantity int   `json:"quantity"`
}

func (i *ItemQuantity) Scan(src interface{}) error {
	str, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", src)
	}
	var val ItemQuantity
	json.NewDecoder(bytes.NewBuffer(str)).Decode(&val)
	*i = val
	return nil
}

type Order struct {
	Id        int64          `json:"id"`
	MessageId int64          `json:"messageId"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Items     []ItemQuantity `json:"items"`
	StateId   OrderStateId   `json:"stateId"`
	Client    User           `json:"client"`
}

type PostOrderDto struct {
	ClientId       int64
	Items          []ItemQuantity `json:"items"`
	ConversationId int64
}

type PatchOrderDto struct {
	Items   []ItemQuantity `json:"items,omitempty"`
	StateId OrderStateId   `json:"stateId,omitempty"`
}

type OrderStateId int

const (
	OrderStateCreated              OrderStateId = 1
	OrderStateAccepted             OrderStateId = 2
	OrderStateDeclined             OrderStateId = 3
	OrderStateFulfilled            OrderStateId = 4
	OrderStateConfirmedFulfillment OrderStateId = 5
	OrderStateSupplierChanges      OrderStateId = 6
	OrderStateClientChanges        OrderStateId = 7
)

var OrderStateMessage = map[OrderStateId]string{
	OrderStateCreated:   "created",
	OrderStateAccepted:  "accepted",
	OrderStateDeclined:  "declined",
	OrderStateFulfilled: "fulfilled",
}

func ValidatePostOrderInput(v *validator.Validator, dto PostOrderDto) {
	v.Check(len(dto.Items) > 0, "itemIds", "must have at least 1 item")
	v.Check(validateQuantity(dto.Items), "itemIds", "quantity must be > 0")
}

func ValidatePatchOrderInput(v *validator.Validator, dto PatchOrderDto) {
	hasItems := len(dto.Items) > 0
	validQuantity := validateQuantity(dto.Items)
	validState := dto.StateId > 0 && dto.StateId <= 7
	if !validQuantity {
		v.AddError("itemIds", "quantity must be > 0")
		return
	}
	if hasItems && validState {
		v.AddError("itemIds", "can't change both items and state")
		v.AddError("stateId", "can't change both items and state")
	}
	if !hasItems && !validState {
		v.AddError("itemIds", "valid items or state change is required")
		v.AddError("stateId", "valid items or state change is required")
	}
}

func validateQuantity(iq []ItemQuantity) bool {
	for _, item := range iq {
		if item.Quantity <= 0 {
			return false
		}
	}
	return true
}
