package app

import "encoding/json"

const (
	EventMessage        = "message"
	EventNewOrder       = "new_order"
	EventUpdateOrder    = "update_order"
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

type WsMessage struct {
	Sender  *Client
	Payload json.RawMessage
}
