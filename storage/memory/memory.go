package memory

import (
	"reflect"
	"github.com/leeif/mercury/storage/data"
	avl "github.com/Workiva/go-datastructures/tree/avl"
)

type MemberInMemory struct {
	Member *avl.Immutable
}

func (m *MemberInMemory) Insert(members ...interface{}) {
	entries := make([]avl.Entry, len(members))
	for i := range entries {
		entries[i] = members[i].(avl.Entry)
	}
	m.Member, _ = m.Member.Insert(entries...)
}

func (m *MemberInMemory) Get(mid ...string) []avl.Entry {
	entries := make([]avl.Entry, len(mid))
	for i := range entries {
		entries[i] = &data.MemberBase{ID: mid[i]}
	}
	res := m.Member.Get(entries...)
	return res
}

type RoomInMemory struct {
	Room *avl.Immutable
}

func (r *RoomInMemory) Insert(rooms ...interface{}) {
	entries := make([]avl.Entry, len(rooms))
	for i := range entries {
		entries[i] = rooms[i].(avl.Entry)
	}
	r.Room, _ = r.Room.Insert(entries...)
}

func (r *RoomInMemory) Get(rid ...string) []avl.Entry {
	entries := make([]avl.Entry, len(rid))
	for i := range entries {
		entries[i] = &data.RoomBase{ID: rid[i]}
	}
	res := r.Room.Get(entries...)
	return res
}

type MessageInMemory struct {
	Message map[string][]interface{}
}

func (m *MessageInMemory) Insert(rid string, message ...interface{}) {
	if m.Message[rid] == nil {
		m.Message[rid] = make([]interface{}, 0)
	}
	m.Message[rid] = append(m.Message[rid], message...)
}

func (m *MessageInMemory) GetHistory(rid string, msg_id int, offset int) []interface{} {
	var end int
	for i, msg := range m.Message[rid] {
		v := reflect.ValueOf(msg)
		if int(reflect.Indirect(v).FieldByName("ID").Int()) == msg_id {
			end = i
			break
		}
	}
	if end-offset < 0 {
		return m.Message[rid][:end]
	}
	return m.Message[rid][(end - offset - 1):end]
}

func (m *MessageInMemory) GetUnRead(rid string, position int) []interface{} {
	if m.Message[rid] != nil && len(m.Message[rid]) > position {
		return m.Message[rid][position:]
	}
	return nil
}

type TokenInMemory struct {
	Token map[string]string
}

func (t *TokenInMemory) Insert(token string, id string) {
	t.Token[token] = id
}

func (t *TokenInMemory) Get(token string) string {
	return t.Token[token]
}

type IndexInMemory struct {
	RoomMember        map[string]map[string]bool
	MemberRoom        map[string]map[string]bool
	RommMemberMessage map[string]int
}

func (i *IndexInMemory) SetMemberOfRoom(rid string, mid string) {
	if i.RoomMember[rid] == nil {
		i.RoomMember[rid] = make(map[string]bool)
	}
	i.RoomMember[rid][mid] = true
}

func (i *IndexInMemory) GetMemberFromRoom(rid string) []string {
	res := make([]string, 0)
	for k := range i.RoomMember[rid] {
		res = append(res, k)
	}
	return res
}

func (i *IndexInMemory) SetRoomOfMember(mid string, rid string) {
	if i.MemberRoom[mid] == nil {
		i.MemberRoom[mid] = make(map[string]bool)
	}
	i.MemberRoom[mid][rid] = true
}

func (i *IndexInMemory) GetRoomFromMember(mid string) []string {
	res := make([]string, 0)
	for k := range i.MemberRoom[mid] {
		res = append(res, k)
	}
	return res
}

func (i *IndexInMemory) SetRoomMemberMessage(rid string, mid string, position int) {
	rmid := rid + ":" + mid
	i.RommMemberMessage[rmid] = position
}

func (i *IndexInMemory) GetRoomMemberMessage(rid string, mid string) int {
	rmid := rid + ":" + mid
	return i.RommMemberMessage[rmid]
}
