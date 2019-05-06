package house

import (
	"mercury/storage"
	"mercury/utils"
)

var house House

func init() {
	house = House{}
	house.store = storage.NewStore()
}

type House struct {
	store  *storage.Store
}

func (house *House) RoomAdd(roomID string, members []string) {
	room := house.GetRoom(roomID)
	// add members' ids into room
	for _, v := range members {
		member := house.GetMember(v)
		house.store.Index.SetMemberOfRoom(room.ID, member.ID)
		house.store.Index.SetRoomOfMember(member.ID, room.ID)
	}
}

func (house *House) RoomDelete(roomID string) {

}

func (house *House) roomMessage(message *Message) {
	for _, roomID := range message.RID {
		room := house.GetRoom(roomID)
		if room != nil {
			room.receiveMessage <- message
		}
	}
}

func (house *House) GetRoom(id string) *Room {
	res := house.store.Room.Get(id)
	if len(res) > 0 && res[0] != nil {
		utils.Debug("room : %v", res[0].(*Room).ID)
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
	house.store.Room.Insert(newRoom)

	return newRoom
}

func (house *House) GetMember(id string) *Member {
	res := house.store.Member.Get(id)
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
	house.store.Member.Insert(newMember)
	return newMember
}

func (house *House) GetMemberFromToken(token string) *Member {
	id := house.store.Token.Get(token)

	if id == "" {
		return nil
	}

	res := house.store.Member.Get(id)
	utils.Debug("%v", res[0])
	if len(res) > 0 && res[0] != nil {
		return res[0].(*Member)
	}
	return nil
}

func (house *House) GetToken(id string) string {
	token := utils.TokenGenerator(id)
	house.store.Token.Insert(token, id)
	return token
}

func GetHouse() *House {
	return &house
}
