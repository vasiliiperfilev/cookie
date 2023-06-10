package app

import (
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/vasiliiperfilev/cookie/internal/data"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

	client := &Client{User: *user, conn: conn, hub: hub, messages: make(chan data.Message, 256)}
	client.hub.register <- client

	go client.readPump()
	go client.writePump()
}
