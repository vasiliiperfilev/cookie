package app

import (
	"github.com/vasiliiperfilev/cookie/internal/data"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan data.Message
	register   chan *Client
	unregister chan *Client
	app        *Application
}

func newHub(app *Application) *Hub {
	return &Hub{
		broadcast:  make(chan data.Message, 256),
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
		case message := <-h.broadcast:
			err := h.app.models.Message.Insert(message)
			if err != nil {
				// TODO: create error message
				// send back error message
				continue
			}
			for client := range h.clients {
				if client.User.Id != message.SenderId {
					client.messages <- message
				}
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				client.conn.Close()
				delete(h.clients, client)
			}
		}
	}
}
