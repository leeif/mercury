package connection

import (
	"net/http"
	"sync"

	"github.com/leeif/mercury/config"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/websocket"
)

// Message type
const (
	// REGISTERROOM : register user into a chat room
	MSG = iota
	CLOSE
)

// Connection struct for each websocket connection
type Connection struct {
	ws        *websocket.Conn
	logger    log.Logger
	Cid       int
	sendChan  chan []byte
	closeChan chan struct{}
	once      sync.Once
	Closed    bool
}

// Reader :
func (c *Connection) Reader(callback func(int, []byte)) {
	defer func() {
		c.close()
		callback(CLOSE, nil)
	}()
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			level.Warn(c.logger).Log("readError", err)
			close(c.closeChan)
			return
		}
		level.Debug(c.logger).Log("recvMsg", message)
		callback(MSG, message)
	}
}

// Writer :
func (c *Connection) Writer(callback func(int, []byte)) {
	defer func() {
		c.close()
	}()
	for {
		select {
		case message, ok := <-c.sendChan:
			if !ok {
				// Close the connection
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				level.Warn(c.logger).Log("msg", "writer close")
				return
			}

			w, err := c.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				level.Debug(c.logger).Log("closeError", err)
				return
			}
		}
	}
}

func (c *Connection) close() {
	c.once.Do(func() {
		close(c.sendChan)
		c.ws.Close()
		c.Closed = true
	})
}

func (c *Connection) SendMessage(message []byte) {
	c.sendChan <- message
}

type Pool struct {
	cid    int
	mutex  sync.Mutex
	logger log.Logger
	config *config.Config
	cons   []*Connection
}

func (p *Pool) New(w http.ResponseWriter, r *http.Request) (*Connection, error) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		level.Error(p.logger).Log("upgradeError", err)
		return nil, err
	}

	conn := &Connection{
		ws:        ws,
		sendChan:  make(chan []byte, 256),
		closeChan: make(chan struct{}),
		logger:    p.logger,
	}

	p.increCID(conn)
	p.cons = append(p.cons, conn)
	return conn, nil
}

func (p *Pool) increCID(conn *Connection) {
	p.mutex.Lock()
	p.cid++
	conn.Cid = p.cid
	p.mutex.Unlock()
}

func NewPool(config *config.Config, l log.Logger) *Pool {
	return &Pool{
		config: config,
		logger: log.With(l, "component", "connection"),
	}
}
