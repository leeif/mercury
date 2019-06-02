package house

import (
	"net/http"

	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/websocket"
	c "github.com/leeif/mercury/connection"
	"github.com/leeif/mercury/storage/data"
)

type Member struct {
	data.MemberBase
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
	t, item, err := newMessage(data)
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
		level.Debug(logger).Log("offset", history.Offest)
		messages := house.roomHistory(history)
		for _, message := range messages {
			msg := &Message{MessageBase: *message}
			if b, err := msg.json(); err == nil {
				member.conn.Send <- b
			}
		}
	}
}

func (member *Member) connClose() {
	member.isClosed = true
}

// GenerateConnection is to handle each websocket connection
func (member *Member) GenerateConnection(w http.ResponseWriter, r *http.Request, connPool *c.Pool) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		level.Error(logger).Log("upgradeError", err)
		return
	}
	member.conn = connPool.New(ws)
	member.isClosed = false
	go member.conn.Reader(member.connCallback)
	go member.conn.Writer(member.connCallback)
	rids := house.Store.GetRoomFromMember(member.ID)
	entries := house.Store.GetRoom(rids...)
	for _, v := range entries {
		if v != nil {
			room := v.(*Room)
			room.receivceMember <- member
		}
	}
}
