package app

import (
	"bytes"
	"encoding/json"

	"github.com/vasiliiperfilev/cookie/internal/data"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan WsMessage
	errors     chan WsMessage
	register   chan *Client
	unregister chan *Client
	app        *Application
}

func newHub(app *Application) *Hub {
	return &Hub{
		broadcast:  make(chan WsMessage, 256),
		errors:     make(chan WsMessage, 256),
		register:   make(chan *Client, 256),
		clients:    make(map[*Client]bool),
		unregister: make(chan *Client, 256),
		app:        app,
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case wsMessage := <-h.broadcast:
			event, err := h.readEvent(wsMessage)
			if err != nil {
				h.errors <- h.createErrorMessage(wsMessage.Sender, "Invalid payload")
				continue
			}
			switch event.Type {
			case EventMessage:
				h.handleMessageEvent(event, wsMessage)
			default:
				h.app.logger.Printf("Unsupported websocket event %v", event)
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				close(client.messages)
				delete(h.clients, client)
			}
		case wsMessage := <-h.errors:
			h.app.logger.Printf("Websocker error event %s", wsMessage.Payload)
			wsMessage.Sender.messages <- wsMessage.Payload
		}
	}
}

func (h *Hub) handleMessageEvent(event *WsEvent, wsMessage WsMessage) {
	var message data.Message
	err := readJson(bytes.NewReader(event.Payload), &message)
	if err != nil {
		h.errors <- h.createErrorMessage(wsMessage.Sender, "Invalid payload")
		return
	}
	err = h.app.models.Message.Insert(message)
	if err != nil {
		h.errors <- h.createErrorMessage(wsMessage.Sender, "Server error")
		return
	}
	for client := range h.clients {
		if client != wsMessage.Sender {
			client.messages <- wsMessage.Payload
		}
	}
}

func (h *Hub) readEvent(wsMessage WsMessage) (*WsEvent, error) {
	var event WsEvent
	err := readJson(bytes.NewReader(wsMessage.Payload), &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (h *Hub) createErrorMessage(client *Client, message string) WsMessage {
	data := ErrorResponse{Message: message, Errors: map[string]string{}}
	payload, _ := json.Marshal(data)
	errEvt := WsEvent{Type: EventError, Payload: payload}
	js, _ := json.Marshal(errEvt)
	wsMessage := WsMessage{Sender: client, Payload: js}
	return wsMessage
}
