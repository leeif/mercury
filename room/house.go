package room

var house = House{Rooms: make([]*Room, 0)}

type House struct {
	Rooms []*Room
}

func (house *House) RegisterRoom(message *Message, member *Member) {
	for _, roomID := range message.RoomID {
		room := house.GetRoom(roomID)
		if room == nil {
			room = newRoom(roomID)
			house.Rooms = append(house.Rooms, room)
		}
		room.Members = append(room.Members, member)
		room.receivceMember <- member
	}
}

func (house *House) GetRoom(roomID string) *Room {
	for _, room := range house.Rooms {
		if room.RoomID == roomID {
			return room
		}
	}
	return nil
}

func (house *House) SendToRoom(message *Message) {
	for _, roomID := range message.RoomID {
		room := house.GetRoom(roomID)
		if room != nil {
			room.receiveMessage <- message
		}
	}
}
