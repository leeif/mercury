package house

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/leeif/mercury/storage/data"
)

type Room struct {
	data.RoomBase
	logger              log.Logger
	storage             data.Store
	memberConnectedChan chan *Member
	houseMessageChan    chan *Message
	houseHistoryChan    chan *History
}

func (room *Room) waitMessageReceive() {
	for {
		select {
		case message, ok := <-room.houseMessageChan:
			if ok && message.RID == room.ID {
				room.transferMessage(message)
			}
		}
	}
}

func (room *Room) waitHistoryReceive() {
	for {
		select {
		case history, ok := <-room.houseHistoryChan:
			if ok && history.RID == room.ID {
				room.transferHistory(history.MID, history.MsgID, history.Offest)
			}
		}
	}
}

func (room *Room) waitMemberConnection() {
	for {
		select {
		case member, ok := <-room.memberConnectedChan:
			if ok {
				room.transferUnReadMessage(member)
			}
		}
	}
}

func (room *Room) transferMessage(message *Message) {
	msgID := room.storage.InsertMessage(&message.MessageBase)
	message.ID = msgID
	mid := room.storage.GetMemberFromRoom(room.ID)
	entries := room.storage.GetMember(mid...)
	members := make([]*Member, len(entries))
	for i := range entries {
		members[i] = entries[i].(*Member)
	}
	for _, member := range members {
		if member != nil {
			if !member.conn.Closed {
				if b, err := message.json(); err == nil {
					member.conn.SendMessage(b)
					room.storage.SetRoomMemberMessage(room.ID, member.ID, message.ID)
				}
			}
		}
	}
}

func (room *Room) transferUnReadMessage(member *Member) {
	msgID := room.storage.GetRoomMemberMessage(room.ID, member.ID)
	messages := room.storage.GetUnReadMessage(room.ID, msgID)
	if messages == nil {
		return
	}
	room.send(messages, member)
}

func (room *Room) transferHistory(mid string, msgID int, offset int) {
	messages := room.storage.GetHistoryMessage(room.ID, msgID, offset)
	entries := room.storage.GetMember(mid)
	if len(entries) == 0 || entries[0] == nil {
		level.Debug(room.logger).Log("msg", "no such member")
		return
	}
	member := entries[0].(*Member)
	room.send(messages, member)
}

func (room *Room) send(messages []*data.MessageBase, member *Member) {
	for _, v := range messages {
		message := &Message{MessageBase: *v}
		if b, err := message.json(); err == nil && !member.conn.Closed {
			member.conn.SendMessage(b)
			room.storage.SetRoomMemberMessage(room.ID, member.ID, v.ID)
		}
	}
}

func (room *Room) Work() {
	go room.waitMessageReceive()
	go room.waitMemberConnection()
}
