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
	cid    int
	logger log.Logger
)

type ConnConfig struct {

}

// Connection struct for each websocket connection
type Connection struct {
	Ws       *websocket.Conn
	Cid      int
	Send     chan []byte
	Closed   bool
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

type Pool struct {
	config *ConnConfig
	cons   []*Connection
}

func (p *Pool) New(ws *websocket.Conn) *Connection {
	cid++
	conn := &Connection{Ws: ws, Cid: cid, Send: make(chan []byte, 256)}
	conn.Closed = false
	return conn
}

func NewPool(config *ConnConfig, l log.Logger) *Pool {
	logger = log.With(l, "component", "connection")
	return &Pool{}
}