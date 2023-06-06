package app

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/vasiliiperfilev/cookie/internal/data"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (a *Application) chatWebSocket(w http.ResponseWriter, r *http.Request) {

	conn, _ := wsUpgrader.Upgrade(w, r, nil)
	_, msg, _ := conn.ReadMessage()
	var message data.Message
	json.NewDecoder(bytes.NewReader(msg)).Decode(&message)
	a.models.Message.Insert(message)
}
