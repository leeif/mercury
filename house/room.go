package house

import (
	"mercury/utils"
	"mercury/storage"
)

type Room struct {
	storage.RoomBase
	receivceMember chan *Member
	receiveMessage chan *Message
}

func (room *Room) ReceiveMessage() {
	for {
		select {
		case message, ok := <-room.receiveMessage:
			if ok {
				room.TransferMessage(message)
			}
		}
	}
}

func (room *Room) ReceiveMember() {
	for {
		select {
		case member, ok := <-room.receivceMember:
			if ok {
				room.TransferUnReadMessage(member)
			}
		}
	}
}

func (room *Room) TransferMessage(message *Message) {
	mid := house.store.Index.GetMemberFromRoom(room.ID)
	entries := house.store.Member.Get(mid...)
	members := make([]*Member, len(entries))
	for i := range entries {
		members[i] = entries[i].(*Member)
	}
	for _, member := range members {
		if member != nil {
			if !member.isClosed {
				utils.Debug("member id : %s", member.ID)
				member.conn.Send <- []byte(message.Text)
				// position increament, should be locked in the furture
			}
		}
	}
}

func (room *Room) TransferUnReadMessage(member *Member) {

}

func (room *Room) Work() {
	go room.ReceiveMessage()
	go room.ReceiveMember()
}

