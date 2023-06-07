package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/vasiliiperfilev/cookie/internal/data"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	UserId int64
	Hub    *Hub
	Conn   *websocket.Conn
}

func (a *Application) chatWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	user, err := a.AuthenticateRequest(w, r)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnathorized):
			a.invalidAuthenticationTokenResponse(w, r)
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	conn, _ := wsUpgrader.Upgrade(w, r, nil)

	addClient(conn, user, hub)

	for {
		message, err := readSentMessage(conn, user)
		if err != nil {
			a.serverErrorResponse(w, r, err)
			return
		}

		err = a.models.Message.Insert(message)
		if err != nil {
			a.serverErrorResponse(w, r, err)
			return
		}

		hub.broadcast <- message
	}
}

func readSentMessage(conn *websocket.Conn, user *data.User) (data.Message, error) {
	_, msg, err := conn.ReadMessage()
	var message data.Message
	json.NewDecoder(bytes.NewReader(msg)).Decode(&message)
	message.SenderId = user.Id
	return message, err
}

func addClient(conn *websocket.Conn, user *data.User, hub *Hub) {
	client := Client{UserId: user.Id, Conn: conn, Hub: hub}
	client.Hub.register <- &client
}

type Hub struct {
	clients   map[int64]Client
	broadcast chan data.Message
	register  chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast: make(chan data.Message),
		register:  make(chan *Client),
		clients:   make(map[int64]Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.UserId] = *client
		case message := <-h.broadcast:
			js, err := json.Marshal(message)
			if err != nil {
				return
			}
			for _, client := range h.clients {
				if client.UserId != message.SenderId {
					client.Conn.WriteMessage(websocket.TextMessage, js)
				}
			}
		}
	}
}
