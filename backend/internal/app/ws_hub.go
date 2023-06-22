package app

import (
	"bytes"
	"encoding/json"

	"github.com/vasiliiperfilev/cookie/internal/data"
	"golang.org/x/exp/slices"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan WsEvent
	errors     chan WsEvent
	register   chan *Client
	unregister chan *Client
	app        *Application
}

func newHub(app *Application) *Hub {
	return &Hub{
		broadcast:  make(chan WsEvent, 256),
		errors:     make(chan WsEvent, 256),
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
		case event := <-h.broadcast:
			switch event.Type {
			case EventMessage:
				h.handleMessageEvent(event)
			default:
				h.app.logger.Printf("Unsupported websocket event %v, payload %v", event.Type, string(event.Payload))
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				close(client.messages)
				delete(h.clients, client)
			}
		case event := <-h.errors:
			h.app.logger.Printf("Websocker error event %s", string(event.Payload))
			event.Sender.messages <- event
		}
	}
}

func (h *Hub) handleMessageEvent(event WsEvent) {
	var dto data.PostMessageDto
	err := readJson(bytes.NewReader(event.Payload), &dto)
	if err != nil {
		h.errors <- h.createErrorMessage(event.Sender, PayloadErrorMessage)
		return
	}
	msg := data.Message{
		Content:        dto.Content,
		ConversationId: dto.ConversationId,
		PrevMessageId:  dto.PrevMessageId,
		SenderId:       event.Sender.User.Id,
	}
	err = h.app.models.Message.Insert(&msg)
	if err != nil {
		h.errors <- h.createErrorMessage(event.Sender, ServerErrorMessage)
		return
	}
	payload, _ := json.Marshal(msg)
	msgEvt := WsEvent{
		Type:    EventMessage,
		Payload: payload,
	}
	var conversation data.Conversation
	if c, ok := event.Sender.Conversations[dto.ConversationId]; !ok {
		conv, err := h.app.models.Conversation.GetById(msg.ConversationId)
		if err != nil {
			h.errors <- h.createErrorMessage(event.Sender, PayloadErrorMessage)
			return
		}
		event.Sender.Conversations[dto.ConversationId] = conv
		conversation = conv
	} else {
		conversation = c
	}

	for client := range h.clients {
		if client != event.Sender && slices.Contains(conversation.UserIds, client.User.Id) {
			client.Conversations[dto.ConversationId] = conversation
			client.messages <- msgEvt
		}
	}
}

func (h *Hub) createErrorMessage(client *Client, message string) WsEvent {
	data := ErrorResponse{Message: message, Errors: map[string]string{}}
	payload, _ := json.Marshal(data)
	errEvt := WsEvent{Type: EventError, Payload: payload, Sender: client}
	return errEvt
}
