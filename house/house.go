package house

import (
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	c "github.com/leeif/mercury/connection"
	"github.com/leeif/mercury/storage/data"
)

var (
	house *House
)

type House struct {
	logger      log.Logger
	Storage     data.Store
	ConnPool    *c.Pool
	messageChan chan *Message
	historyChan chan *History
}

func (house *House) RoomAdd(roomID string, members []string) {
	room := house.GetRoom(roomID)
	// add members' ids into room
	for _, v := range members {
		member := house.GetMember(v)
		// insert new room into storage engin
		house.Storage.InsertRoomMember(room, member)
		house.Storage.SetMemberOfRoom(room.ID, member.ID)
		house.Storage.SetRoomOfMember(member.ID, room.ID)
	}
}

func (house *House) RoomDelete(roomID string) {

}

func (house *House) GetRoom(id string) *Room {
	res := house.Storage.GetRoom(id)
	if len(res) > 0 && res[0] != nil {
		level.Debug(house.logger).Log("roomID", res[0].(*Room).ID)
		return res[0].(*Room)
	}
	return house.NewRoom(id)
}

func (house *House) NewRoom(id string) *Room {
	newRoom := &Room{
		storage:             house.Storage,
		houseMessageChan:    house.messageChan,
		memberConnectedChan: make(chan *Member, 5),
	}
	newRoom.ID = id
	newRoom.Work()

	house.Storage.InsertRoom(newRoom)
	return newRoom
}

func (house *House) GetMember(id string) *Member {
	res := house.Storage.GetMember(id)
	if len(res) > 0 && res[0] != nil {
		return res[0].(*Member)
	}

	// notify all the rooms that this member has joined
	return house.NewMember(id)
	// return nil
}

func (house *House) NewMember(id string) *Member {
	newMember := &Member{
		storage:          house.Storage,
		houseHistoryChan: house.historyChan,
		houseMessageChan: house.messageChan,
	}
	newMember.ID = id
	newMember.logger = log.With(house.logger, "component", "member")
	house.Storage.InsertMember(newMember)
	return newMember
}

func (house *House) MemberConnect(w http.ResponseWriter, r *http.Request, mid string, token string) error {
	member := house.GetMember(mid)
	if !member.verifyToken(token) {
		return errors.New("invalid token")
	}
	err := member.GenerateConnection(w, r, house.ConnPool)
	if err != nil {
		return err
	}
	return nil
}

func (house *House) NewToken(mid string) string {
	member := house.GetMember(mid)
	return member.newToken()
}

func NewHouse(l log.Logger, store data.Store, connPool *c.Pool) *House {
	if house == nil {
		house = &House{
			Storage:     store,
			ConnPool:    connPool,
			logger:      log.With(l, "component", "house"),
			messageChan: make(chan *Message, 100),
			historyChan: make(chan *History, 100),
		}
	}
	return house
}
