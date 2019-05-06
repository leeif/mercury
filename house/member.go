package house

import (
	"log"
	c "mercury/connection"
	"mercury/storage"
	"mercury/utils"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	// global connection id
	cid = 0
)

type Member struct {
	storage.MemberBase
	conn     *c.Connection
	isClosed bool
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
		return
	}
	message.MID = member.ID
	house.roomMessage(message)
}

func (member *Member) connClose() {
	utils.Debug("connection is closed, member id : %s", member.ID)
	member.isClosed = true
}

// GenerateConnection is to handle each websocket connection
func (member *Member) GenerateConnection(w http.ResponseWriter, r *http.Request) {
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
	member.isClosed = false
	go member.conn.Reader(member.connCallback)
	go member.conn.Writer(member.connCallback)
	rids := house.store.Index.GetRoomFromMember(member.ID)
	entries := house.store.Room.Get(rids...)
	for _, v := range entries {
		if v != nil {
			room := v.(*Room)
			room.receivceMember <- member
		}
	}
}
