package room

import (
	"log"
	c "mercury/connection"
	"net/http"
	"mercury/utils"
	"github.com/gorilla/websocket"
)

var (
	// global connection id
	cid = 0
)

type Member struct {
	id           string
	registerTime string
	readPosition string
	conn         *c.Connection
	isClosed     bool
}

func (member *Member) connCallback(flag int, data []byte) {
	switch flag {
	case c.FLAG_RECEV_MESSAGE:
		member.connRecevMessage(data)
	case c.FLAG_CONN_CLOSE:
		member.connClose()
	}
}

func (member *Member) connRecevMessage(data []byte) {
	message, err := NewMessage(data)
	if err != nil {
		member.sendResponse()
		return
	}
	message.FromID = member.id
	switch message.MsgType {
	case REGISTERROOM:
		member.registerRoom(message)
	case SENDMESSAGE:
		member.sendMessage(message)
	}
}

func (member *Member) connClose() {
	utils.Debug("connection is closed, member id : %s", member.id)
	member.isClosed = true
}

func (member *Member) registerRoom(message *Message) {
	house.RegisterRoom(message, member)
}

func (member *Member) sendMessage(message *Message) {
	house.SendToRoom(message)
}

func (member *Member) sendResponse() {


}

// GenerateConnection is to handle each websocket connection
func (member Member) GenerateConnection(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	cid++
	member.conn = &c.Connection{Ws: ws, Cid: cid, Send: make(chan []byte, 256)}
	member.id = r.URL.Path
	go member.conn.Reader(member.connCallback)
	go member.conn.Writer(member.connCallback)
}

func NewMember() *Member {
	member := &Member{}
	return member
}