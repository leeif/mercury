package house

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/leeif/mercury/common"
	c "github.com/leeif/mercury/connection"
	"github.com/leeif/mercury/storage/data"
)

var (
	house  *House
	logger log.Logger
)

type House struct {
	Store    data.Store
	ConnPool *c.Pool
}

func (house *House) RoomAdd(roomID string, members []string) {
	room := house.GetRoom(roomID)
	// add members' ids into room
	for _, v := range members {
		member := house.GetMember(v)
		// insert new room into storage engin
		house.Store.InsertRoomMember(room, member)
		house.Store.SetMemberOfRoom(room.ID, member.ID)
		house.Store.SetRoomOfMember(member.ID, room.ID)
	}
}

func (house *House) RoomDelete(roomID string) {

}

func (house *House) roomMessage(message *Message) {
	room := house.GetRoom(message.RID)
	if room != nil {
		msg_id := house.Store.InsertMessage(&message.MessageBase)
		message.ID = msg_id
		room.receiveMessage <- message
	}
}

func (house *House) roomHistory(history *History) []*data.MessageBase {
	messages := house.Store.GetHistoryMessage(history.RID, history.MsgID, history.Offest)
	return messages
}

func (house *House) GetRoom(id string) *Room {
	res := house.Store.GetRoom(id)
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

	house.Store.InsertRoom(newRoom)
	return newRoom
}

func (house *House) GetMember(id string) *Member {
	res := house.Store.GetMember(id)
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
	house.Store.InsertMember(newMember)
	return newMember
}

func (house *House) GetMemberFromToken(token string) *Member {
	id := house.Store.GetToken(token)

	if id == "" {
		return nil
	}

	res := house.Store.GetMember(id)
	if len(res) > 0 && res[0] != nil {
		return res[0].(*Member)
	}
	return house.NewMember(id)
}

func (house *House) GetToken(id string) string {
	token := common.TokenGenerator(id)
	house.Store.InsertToken(token, id)
	level.Debug(logger).Log("token", token)
	return token
}

func NewHouse(l log.Logger, s data.Store, connPool *c.Pool) *House {
	if house == nil {
		house = &House{
			Store: s,
		}
	}
	logger = log.With(l, "component", "house")
	return house
}
