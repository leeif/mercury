package house

import (
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	c "github.com/leeif/mercury/connection"
	"github.com/leeif/mercury/district/child"
	"github.com/leeif/mercury/storage/data"
)

type House struct {
	logger   log.Logger
	storage  data.Store
	connPool *c.Pool
}

func (house *House) RoomAdd(rid string, members []string) {
	for _, mid := range members {
		// insert new room into storage engin
		house.storage.InsertRoomMember(rid, mid)
		house.storage.SetMemberOfRoom(rid, mid)
		house.storage.SetRoomOfMember(mid, rid)
	}
}

func (house *House) RoomDelete(roomID string) {

}

func (house *House) getMember(id string) *Member {
	res := house.storage.GetMember(id)
	if len(res) > 0 && res[0] != nil {
		return res[0].(*Member)
	}

	return house.newMember(id)
}

func (house *House) newMember(id string) *Member {
	newMember := &Member{
		storage: house.storage,
	}
	newMember.ID = id
	newMember.logger = log.With(house.logger, "component", "member")
	house.storage.InsertMember(newMember)
	return newMember
}

func (house *House) MemberConnect(w http.ResponseWriter, r *http.Request, mid string, token string) error {
	member := house.getMember(mid)
	if !member.verifyToken(token) {
		return errors.New("invalid token")
	}
	err := member.GenerateConnection(w, r, house.connPool)
	if err != nil {
		return err
	}

	rids := house.storage.GetRoomFromMember(member.ID)
	for _, rid := range rids {
		room := newRoom(rid, house.storage)
		room.transferUnReadMessage(member)
	}

	child.RegisterMember()

	go member.conn.Reader(member.connCallback)
	go member.conn.Writer(member.connCallback)

	return nil
}

func (house *House) NewToken(mid string) string {
	member := house.getMember(mid)
	return member.newToken()
}

func NewHouse(l log.Logger, store data.Store, connPool *c.Pool) *House {
	house := &House{
		storage:  store,
		connPool: connPool,
		logger:   log.With(l, "component", "house"),
	}
	return house
}
