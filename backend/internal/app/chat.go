package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vasiliiperfilev/cookie/internal/data"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	UserId int64
	hub    *Hub
	conn   *websocket.Conn
	// messages chan data.Message
}

// func (c *Client) ReadPump() {
// 	defer func() {
// 		c.conn.Close()
// 	}()
// 	c.conn.SetReadLimit(maxMessageSize)
// 	c.conn.SetReadDeadline(time.Now().Add(pongWait))
// 	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
// 	select {
// 	case message := <-c.messages:
// 		c.hub.broadcast <- message
// 	default:
// 		c.conn.Close()
// 	}
// }

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
			conn.Close()
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

func addClient(conn *websocket.Conn, user *data.User, hub *Hub) *Client {
	client := &Client{UserId: user.Id, conn: conn, hub: hub}
	client.hub.register <- client
	return client
}

type Hub struct {
	clients   map[int64]Client
	broadcast chan data.Message
	register  chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast: make(chan data.Message, 255),
		register:  make(chan *Client, 255),
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
					err := client.conn.WriteMessage(websocket.TextMessage, js)
					if err != nil {
						client.conn.Close()

					}
				}
			}
		}
	}
}
