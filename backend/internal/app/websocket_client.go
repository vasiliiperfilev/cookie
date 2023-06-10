package app

import (
	"bytes"
	"encoding/json"
	"log"
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
	messages chan data.Message
}

func (c *Client) readPump() {
	defer func() {
		// Graceful Close the Connection once this
		// function is done
		c.hub.unregister <- c
		c.conn.Close()
	}()
	// Loop Forever
	for {
		// ReadMessage is used to read the next message in queue
		// in the connection
		message, err := readSentMessage(c)
		c.conn.SetReadLimit(maxMessageSize)
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
		if err != nil {
			// If Connection is closed, we will Recieve an error here
			// We only want to log Strange errors, but not simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// TODO: send error to hub to log and send message to client itself
				log.Printf("error reading message: %v", err)
			}
			break // Break the loop to close conn & Cleanup
		}
		c.hub.broadcast <- *message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	// if message comes before the end of ping period work on it
	// if ping period came earlier - send ping to check if client alive
	for {
		select {
		case message, ok := <-c.messages:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			// initialize in handler and run tests
			// then go back to handling db write error
			// add tests for client disconnects
			// client stays connected but silent
			// clients sends incorrect message format
			// client disconnects with error
			js, err := json.Marshal(message)
			if err != nil {
				return
			}
			w.Write(js)

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

func readSentMessage(client *Client) (*data.Message, error) {
	_, msg, err := client.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	var message data.Message
	err = readJson(bytes.NewReader(msg), &message)
	message.SenderId = client.User.Id
	if err != nil {
		return nil, err
	}
	return &message, err
}
