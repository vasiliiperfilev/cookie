package app

import "encoding/json"

const (
	EventMessage        = "message"
	EventError          = "error"
	PayloadErrorMessage = "Invalid payload"
	ServerErrorMessage  = "Server error"
)

// WsEvent is the Messages sent over the websocket
// Used to differ between different actions
type WsEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type WsMessage struct {
	Sender  *Client
	Payload json.RawMessage
}
