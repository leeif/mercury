package house

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/leeif/mercury/common"
	c "github.com/leeif/mercury/connection"
	"github.com/leeif/mercury/storage"
)

var (
	house  *House
	logger log.Logger
)

type House struct {
	Store    *storage.Store
	ConnPool *c.Pool
}

func (house *House) RoomAdd(roomID string, members []string) {
	room := house.GetRoom(roomID)
	// add members' ids into room
	for _, v := range members {
		member := house.GetMember(v)
		house.Store.Index.SetMemberOfRoom(room.ID, member.ID)
		house.Store.Index.SetRoomOfMember(member.ID, room.ID)
	}
}

func (house *House) RoomDelete(roomID string) {

}

func (house *House) roomMessage(message *Message) {
	room := house.GetRoom(message.RID)
	if room != nil {
		room.receiveMessage <- message
		house.Store.Message.Insert(room.ID, message)
	}
}

func (house *House) roomHistory(history *History) []interface{} {
	messages := house.Store.Message.GetHistory(history.RID, history.MsgID, history.Offest)
	return messages
}

func (house *House) GetRoom(id string) *Room {
	res := house.Store.Room.Get(id)
	if len(res) > 0 && res[0] != nil {
		level.Debug(logger).Log("roomID", res[0].(*Room).ID)
		return res[0].(*Room)
	}
	return house.NewRoom(id)
}

func (house *House) NewRoom(id string) *Room {
	newRoom := &Room{
		receiveMessage: make(chan *Message, 10),
		receivceMember: make(chan *Member, 5),
	}
	newRoom.ID = id
	newRoom.Work()
	// insert new room into storage engin
	house.Store.Room.Insert(newRoom)

	return newRoom
}

func (house *House) GetMember(id string) *Member {
	res := house.Store.Member.Get(id)
	if len(res) > 0 && res[0] != nil {
		return res[0].(*Member)
	}

	// notify all the rooms that this member has joined
	return house.NewMember(id)
	// return nil
}

func (house *House) NewMember(id string) *Member {
	newMember := &Member{
		isClosed: true,
	}
	newMember.ID = id
	// insert new member into avl
	house.Store.Member.Insert(newMember)
	return newMember
}

func (house *House) GetMemberFromToken(token string) *Member {
	id := house.Store.Token.Get(token)

	if id == "" {
		return nil
	}

	res := house.Store.Member.Get(id)
	if len(res) > 0 && res[0] != nil {
		return res[0].(*Member)
	}
	return nil
}

func (house *House) GetToken(id string) string {
	token := common.TokenGenerator(id)
	house.Store.Token.Insert(token, id)
	return token
}

func NewHouse(l log.Logger, s *storage.Store, connPool *c.Pool) *House {
	if house == nil {
		house = &House{
			Store: s,
		}
	}
	logger = log.With(l, "component", "house")
	return house
}
