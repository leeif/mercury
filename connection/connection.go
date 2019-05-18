package connection

import (
	"github.com/gorilla/websocket"
	"github.com/leeif/mercury/utils"
)

// Message type
const (
	// REGISTERROOM : register user into a chat room
	MSG = iota
	// SENDTOROOM : send a message to a room
	CLOSE
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
			utils.Info("read:", err)
			break
		}
		utils.Info("recv: %s", message)
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
