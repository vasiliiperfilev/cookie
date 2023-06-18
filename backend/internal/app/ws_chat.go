package app

import (
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/vasiliiperfilev/cookie/internal/data"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (a *Application) wsChatHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	user, err := a.AuthenticateWsUpgradeRequest(w, r)
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
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		a.logger.Print(err)
		return
	}

	client := &Client{User: *user, conn: conn, hub: hub, messages: make(chan WsEvent, 256)}
	client.hub.register <- client

	go client.readPump()
	go client.writePump()
}
