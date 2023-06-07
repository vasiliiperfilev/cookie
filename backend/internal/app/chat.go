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
	Conn *websocket.Conn
}

var clients = map[int64]Client{}

func (a *Application) chatWebSocket(w http.ResponseWriter, r *http.Request) {
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

	client := Client{Conn: conn}
	clients[user.Id] = client

	_, msg, _ := conn.ReadMessage()
	var message data.Message
	json.NewDecoder(bytes.NewReader(msg)).Decode(&message)
	message.SenderId = user.Id
	a.models.Message.Insert(message)
	// get conversation
	// send message to everyone within conversation
	conversation, err := a.models.Conversation.GetById(message.ConversationId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	for _, userId := range conversation.UserIds {
		js, err := json.Marshal(message)
		if err != nil {
			a.notFoundResponse(w, r)
			return
		}
		if client, ok := clients[userId]; ok {
			client.Conn.WriteMessage(websocket.TextMessage, js)
		}
	}
}
