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
			case EventConfirmation:
				h.handleConfirmationEvent(event)
			case EventOrder:
				h.handleOrderEvent(event)
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
	conversation, err := h.getConversation(event, msg)
	if err != nil {
		h.errors <- h.createErrorMessage(event.Sender, PayloadErrorMessage)
		return
	}

	for client := range h.clients {
		if slices.Contains(conversation.UserIds, client.User.Id) {
			client.Conversations[dto.ConversationId] = conversation
			client.messages <- msgEvt
		}
	}
}

func (h *Hub) handleConfirmationEvent(event WsEvent) {
	var wsConfirm struct{ MessageId int64 }
	err := readJson(bytes.NewReader(event.Payload), &wsConfirm)
	if err != nil {
		h.errors <- h.createErrorMessage(event.Sender, PayloadErrorMessage)
		return
	}

}

func (h *Hub) handleOrderEvent(event WsEvent) {
	var order data.Order
	err := readJson(bytes.NewReader(event.Payload), &order)
	if err != nil {
		h.errors <- h.createErrorMessage(event.Sender, PayloadErrorMessage)
		return
	}

	msg, err := h.app.models.Message.GetById(order.MessageId)
	if err != nil {
		h.errors <- h.createErrorMessage(event.Sender, ServerErrorMessage)
		return
	}

	payload, _ := json.Marshal(order)
	orderEvent := WsEvent{
		Type:    EventOrder,
		Payload: payload,
	}

	conversation, err := h.getConversation(event, msg)
	if err != nil {
		h.errors <- h.createErrorMessage(event.Sender, PayloadErrorMessage)
		return
	}

	for client := range h.clients {
		if slices.Contains(conversation.UserIds, client.User.Id) {
			client.Conversations[msg.ConversationId] = conversation
			client.messages <- orderEvent
		}
	}
}

func (h *Hub) getConversation(event WsEvent, msg data.Message) (data.Conversation, error) {
	if conversation, ok := event.Sender.Conversations[msg.ConversationId]; !ok {
		c, err := h.app.models.Conversation.GetById(msg.ConversationId)
		if err != nil {
			return data.Conversation{}, err
		}
		event.Sender.Conversations[msg.ConversationId] = c
		return c, nil
	} else {
		return conversation, nil
	}
}

func (h *Hub) createErrorMessage(client *Client, message string) WsEvent {
	data := ErrorResponse{Message: message, Errors: map[string]string{}}
	payload, _ := json.Marshal(data)
	errEvt := WsEvent{Type: EventError, Payload: payload, Sender: client}
	return errEvt
}
