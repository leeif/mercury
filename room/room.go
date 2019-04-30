package room

import (
	"mercury/utils"
)

type Room struct {
	RoomID         string
	Members        []*Member
	receivceMember chan *Member
	UnReadMessage  []*Message
	receiveMessage chan *Message

	readPosition map[string]int
}

func (room *Room) ReceiveMessage() {
	for {
		select {
		case message, ok := <-room.receiveMessage:
			if !ok {
			}
			room.TransferMessage(message)
		}
	}
}

func (room *Room) ReceiveMember() {
	for {
		select {
		case member, ok := <-room.receivceMember:
			if !ok {
			}
			room.TransferUnReadMessage(member)
		}
	}
}

func (room *Room) TransferMessage(message *Message) {
	room.UnReadMessage = append(room.UnReadMessage, message)
	for _, member := range room.Members {
		if member != nil && !member.isClosed {
			utils.Debug("member id : ", member.id)
			member.conn.Send <- []byte(message.Body)
			// position increament, should be locked in the furture
			room.readPosition[member.id] = room.readPosition[member.id] + 1
		}
	}
}

func (room *Room) TransferUnReadMessage(member *Member) {
	utils.Debug("read position %d", room.readPosition[member.id])
	for i := room.readPosition[member.id]; i < len(room.UnReadMessage); i++ {
		member.conn.Send <- []byte(room.UnReadMessage[i].Body)
		// position increament, should be locked in the furture
		room.readPosition[member.id] = room.readPosition[member.id] + 1
	}
}

func (room *Room) Work() {
	go room.ReceiveMessage()
	go room.ReceiveMember()
}

func newRoom(roomID string) *Room {
	newRoom := &Room{
		RoomID:         roomID,
		Members:        make([]*Member, 0),
		UnReadMessage:  make([]*Message, 0),
		receiveMessage: make(chan *Message, 10),
		receivceMember: make(chan *Member, 5),
		readPosition:   make(map[string]int),
	}
	newRoom.Work()
	return newRoom
}
