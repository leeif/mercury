package memory

import (
	"sync"

	avl "github.com/Workiva/go-datastructures/tree/avl"
	"github.com/leeif/mercury/storage/data"
)

var (
	msg_id int = 0
)

type Memory struct {
	Member            *avl.Immutable
	Room              *avl.Immutable
	Token             map[string]string
	Message           map[string][]*data.MessageBase
	RoomMember        map[string]map[string]bool
	MemberRoom        map[string]map[string]bool
	RommMemberMessage map[string]int
	msg_id_mutex      sync.Mutex
}

func (m *Memory) InsertMember(members ...interface{}) {
	entries := make([]avl.Entry, len(members))
	for i := range entries {
		entries[i] = members[i].(avl.Entry)
	}
	m.Member, _ = m.Member.Insert(entries...)
}

func (m *Memory) GetMember(mid ...string) []interface{} {
	entries := make([]avl.Entry, len(mid))
	for i := range entries {
		entries[i] = &data.MemberBase{ID: mid[i]}
	}
	members := m.Member.Get(entries...)
	res := make([]interface{}, 0)
	for _, v := range members {
		if v != nil {
			res = append(res, v.(interface{}))
		}
	}
	return res
}

func (m *Memory) InsertRoomMember(room interface{}, member interface{}) {
	// do nothing
	return
}

func (m *Memory) InsertRoom(rooms ...interface{}) {
	entries := make([]avl.Entry, len(rooms))
	for i := range entries {
		entries[i] = rooms[i].(avl.Entry)
	}
	m.Room, _ = m.Room.Insert(entries...)
}

func (m *Memory) GetRoom(rid ...string) []interface{} {
	entries := make([]avl.Entry, len(rid))
	for i := range entries {
		entries[i] = &data.RoomBase{ID: rid[i]}
	}
	rooms := m.Room.Get(entries...)
	res := make([]interface{}, 0)
	for _, v := range rooms {
		if v != nil {
			res = append(res, v.(interface{}))
		}
	}
	return res
}

func (m *Memory) InsertMessage(message *data.MessageBase) int {
	rid := message.RID
	if m.Message[rid] == nil {
		m.Message[rid] = make([]*data.MessageBase, 0)
	}
	m.msg_id_mutex.Lock()
	msg_id++
	message.ID = msg_id
	m.msg_id_mutex.Unlock()
	m.Message[rid] = append(m.Message[rid], message)
	return message.ID
}

func (m *Memory) GetHistoryMessage(rid string, msg_id int, offset int) []*data.MessageBase {
	var end int
	for i, msg := range m.Message[rid] {
		if msg.ID == msg_id {
			end = i
			break
		}
	}
	if end-offset < 0 {
		return m.Message[rid][:end]
	}
	return m.Message[rid][(end - offset - 1):end]
}

func (m *Memory) GetUnReadMessage(rid string, msg_id int) []*data.MessageBase {
	position := -1
	for i, msg := range m.Message[rid] {
		if msg.ID == msg_id {
			position = i
			break
		}
	}
	if position != -1 && m.Message[rid] != nil && len(m.Message[rid]) > position {
		return m.Message[rid][(position+1):]
	}
	return nil
}

func (m *Memory) InsertToken(token string, mid string) {
	m.Token[token] = mid
}

func (m *Memory) GetToken(token string) string {
	return m.Token[token]
}

func (m *Memory) SetMemberOfRoom(rid string, mid string) {
	if m.RoomMember[rid] == nil {
		m.RoomMember[rid] = make(map[string]bool)
	}
	m.RoomMember[rid][mid] = true
}

func (m *Memory) GetMemberFromRoom(rid string) []string {
	res := make([]string, 0)
	for k := range m.RoomMember[rid] {
		res = append(res, k)
	}
	return res
}

func (m *Memory) SetRoomOfMember(mid string, rid string) {
	if m.MemberRoom[mid] == nil {
		m.MemberRoom[mid] = make(map[string]bool)
	}
	m.MemberRoom[mid][rid] = true
}

func (m *Memory) GetRoomFromMember(mid string) []string {
	res := make([]string, 0)
	for k := range m.MemberRoom[mid] {
		res = append(res, k)
	}
	return res
}

func (m *Memory) SetRoomMemberMessage(rid string, mid string, msg_id int) {
	rmid := rid + ":" + mid
	m.RommMemberMessage[rmid] = msg_id
}

func (m *Memory) GetRoomMemberMessage(rid string, mid string) int {
	rmid := rid + ":" + mid
	return m.RommMemberMessage[rmid]
}
