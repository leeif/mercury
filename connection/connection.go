package connection

import (
	"github.com/gorilla/websocket"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Message type
const (
	// REGISTERROOM : register user into a chat room
	MSG = iota
	// SENDTOROOM : send a message to a room
	CLOSE
)

var (
	logger log.Logger
)

// Connection struct for each websocket connection
type Connection struct {
	Ws   *websocket.Conn
	Cid  int
	Send chan []byte
}

// Reader :
func (c *Connection) Reader(callback func(int, []byte)) {
	defer func() {
		c.Ws.Close()
		callback(CLOSE, nil)
	}()
	for {
		_, message, err := c.Ws.ReadMessage()
		if err != nil {
			level.Error(logger).Log("readError",err)
			break
		}
		level.Debug(logger).Log("recvMsg", message)
		callback(MSG, message)
	}
}

// Writer :
func (c *Connection) Writer(callback func(int, []byte)) {
	defer func() {
		c.Ws.Close()
		callback(CLOSE, nil)
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// Close the connection
				c.Ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Ws.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func WithLogger(l log.Logger) {
	logger = log.With(l, "component", "connection")
}