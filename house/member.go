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
	case c.MSG:
		member.connRecevMessage(data)
	case c.CLOSE:
		member.connClose()
	}
}

func (member *Member) connRecevMessage(data []byte) {
	t, item, err := New(data)
	if err != nil {
		return
	}
	switch t {
	case SEND:
		message := item.(*Message)
		message.MID = member.ID
		house.roomMessage(message)
	case HISTORY:
		history := item.(*History)
		messages := house.roomHistory(history)
		res := Response{Status:"ok", Body: messages}
		if b, err := res.json(); err == nil {
			member.conn.Send <- b
		}
	}
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
