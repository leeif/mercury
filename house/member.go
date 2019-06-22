package house

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/leeif/mercury/common"
	c "github.com/leeif/mercury/connection"
	"github.com/leeif/mercury/storage/data"
)

type Member struct {
	data.MemberBase
	logger           log.Logger
	storage          data.Store
	conn             *c.Connection
	houseMessageChan chan *Message
	houseHistoryChan chan *History
}

func (member *Member) connCallback(flag int, data []byte) {
	switch flag {
	case c.MSG:
		member.connRecevMessage(data)
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
		member.broadcastMessage(message)
	case HISTORY:
		history := item.(*History)
		level.Debug(member.logger).Log("offset", history.Offest)
		member.broadcastHistory(history)
	}
}

func (member *Member) broadcastMessage(message *Message) {
	member.houseMessageChan <- message
}

func (member *Member) broadcastHistory(history *History) {
	member.houseHistoryChan <- history
}

func (member *Member) newToken() string {
	token := common.TokenGenerator(member.ID)
	member.storage.InsertToken(member.ID, token)
	level.Debug(member.logger).Log("token", token)
	return token
}

func (member *Member) verifyToken(token string) bool {
	t := member.storage.GetToken(member.ID)
	if t == token {
		return true
	}
	return false
}

// GenerateConnection is to handle each websocket connection
func (member *Member) GenerateConnection(w http.ResponseWriter, r *http.Request, connPool *c.Pool) error {
	var err error
	member.conn, err = connPool.New(w, r)
	if err != nil {
		return err
	}
	go member.conn.Reader(member.connCallback)
	go member.conn.Writer(member.connCallback)

	return nil

	// rids := house.Storage.GetRoomFromMember(member.ID)
	// entries := house.Storage.GetRoom(rids...)
	// for _, v := range entries {
	// 	if v != nil {
	// 		room := v.(*Room)
	// 		room.memberConnectedChan <- member
	// 	}
	// }
}
