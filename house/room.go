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
				if b, err := message.json(); err  == nil {
					member.conn.Send <- b
					// position increament, should be locked in the furture
					position := house.store.Index.GetRoomMemberMessage(room.ID, member.ID)
					house.store.Index.SetRoomMemberMessage(room.ID, member.ID, position+1)
				}
			}
		}
	}
}

func (room *Room) TransferUnReadMessage(member *Member) {
	position := house.store.Index.GetRoomMemberMessage(room.ID, member.ID)
	messages := house.store.Message.GetUnRead(room.ID, position)
	if messages == nil {
		return
	}
	for _, v := range messages {
		message := v.(*Message)
		if b, err := message.json(); err == nil && !member.isClosed {
			member.conn.Send <- b
			// position increament, should be locked in the furture
			position := house.store.Index.GetRoomMemberMessage(room.ID, member.ID)
			house.store.Index.SetRoomMemberMessage(room.ID, member.ID, position+1)
		}
	}
}

func (room *Room) Work() {
	go room.ReceiveMessage()
	go room.ReceiveMember()
}

