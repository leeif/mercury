package house

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	c "github.com/leeif/mercury/connection"
	"github.com/leeif/mercury/storage/data"
	"github.com/leeif/mercury/utils"
)

type Member struct {
	data.MemberBase
	logger  log.Logger
	storage data.Store
	conn    *c.Connection
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
		member.sendMessage(message)
	case HISTORY:
		history := item.(*History)
		level.Debug(member.logger).Log("offset", history.Offest)
		member.sendHistory(history)
	}
}

func (member *Member) sendMessage(message *Message) {
	room := newRoom(message.RID, member.storage)
	room.transferMessage(message)
}

func (member *Member) sendHistory(history *History) {
	room := newRoom(history.RID, member.storage)
	room.transferHistory(member, history)
}

func (member *Member) newToken() string {
	token := utils.TokenGenerator(member.ID)
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

	return nil
}
