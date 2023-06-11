package app

import (
	"bytes"
	"fmt"
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

type Client struct {
	User     data.User
	hub      *Hub
	conn     *websocket.Conn
	messages chan []byte
}

func (c *Client) readPump() {
	defer func() {
		// Graceful Close the Connection once this
		// function is done
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	// Loop Forever
	for {
		// readSentEvent is used to read the next event in queue
		// in the connection
		_, evt, err := c.conn.ReadMessage()
		if err != nil {
			// If Connection is closed, we will Recieve an error here
			// We only want to log Strange errors, but not simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.hub.errors <- c.hub.createErrorMessage(c, fmt.Sprintf("error reading message: %v", err))
			}
			break // Break the loop to close conn & Cleanup
		}
		c.hub.broadcast <- WsMessage{Sender: c, Payload: evt}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	// if event comes before the end of ping period work on it
	// if ping period came earlier - send ping to check if client alive
	for {
		select {
		case msg, ok := <-c.messages:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel after readPump was stopped see defer in readPump
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(msg)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func readSentEvent(client *Client) (*WsEvent, error) {
	_, msg, err := client.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	var event WsEvent
	err = readJson(bytes.NewReader(msg), &event)
	if err != nil {
		return nil, err
	}
	return &event, err
}
