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

	addClient(conn, user)

	message, err := readNextMessage(conn, user)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	err = a.models.Message.Insert(message)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	err = sendMessage(a, message)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
}

func sendMessage(a *Application, message data.Message) error {
	conversation, err := a.models.Conversation.GetById(message.ConversationId)
	for _, userId := range conversation.UserIds {
		js, err := json.Marshal(message)
		if err != nil {
			return err
		}
		if client, ok := clients[userId]; ok && userId != message.SenderId {
			client.Conn.WriteMessage(websocket.TextMessage, js)
		}
	}
	return err
}

func readNextMessage(conn *websocket.Conn, user *data.User) (data.Message, error) {
	_, msg, err := conn.ReadMessage()
	var message data.Message
	json.NewDecoder(bytes.NewReader(msg)).Decode(&message)
	message.SenderId = user.Id
	return message, err
}

func addClient(conn *websocket.Conn, user *data.User) {
	client := Client{Conn: conn}
	clients[user.Id] = client
}
