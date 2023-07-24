package app

import "encoding/json"

const (
	EventMessage        = "message"
	EventOrder          = "order"
	EventConfirmation   = "confirmation"
	EventError          = "error"
	PayloadErrorMessage = "Invalid payload"
	ServerErrorMessage  = "Server error"
)

// WsEvent is the Messages sent over the websocket
// Used to differ between different actions
type WsEvent struct {
	Type    string          `json:"type"`
	Sender  *Client         `json:"-"`
	Payload json.RawMessage `json:"payload"`
}
