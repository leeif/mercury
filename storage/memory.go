package storage

import (
	"mercury/utils"
	"reflect"

	avl "github.com/Workiva/go-datastructures/tree/avl"
)

type MemberInMemory struct {
	member *avl.Immutable
}

func (m *MemberInMemory) Insert(members ...interface{}) {
	entries := make([]avl.Entry, len(members))
	for i := range entries {
		entries[i] = members[i].(avl.Entry)
	}
	m.member, _ = m.member.Insert(entries...)
	utils.Debug("member len %v", m.member.Len())
}

func (m *MemberInMemory) Get(mid ...string) []avl.Entry {
	entries := make([]avl.Entry, len(mid))
	for i := range entries {
		entries[i] = &MemberBase{ID: mid[i]}
	}
	res := m.member.Get(entries...)
	return res
}

type RoomInMemory struct {
	room *avl.Immutable
}

func (r *RoomInMemory) Insert(rooms ...interface{}) {
	entries := make([]avl.Entry, len(rooms))
	for i := range entries {
		entries[i] = rooms[i].(avl.Entry)
	}
	r.room, _ = r.room.Insert(entries...)
}

func (r *RoomInMemory) Get(rid ...string) []avl.Entry {
	entries := make([]avl.Entry, len(rid))
	for i := range entries {
		entries[i] = &RoomBase{ID: rid[i]}
	}
	res := r.room.Get(entries...)
	return res
}

type MessageInMemory struct {
	message map[string][]interface{}
}

func (m *MessageInMemory) Insert(rid string, message ...interface{}) {
	if m.message[rid] == nil {
		m.message[rid] = make([]interface{}, 0)
	}
	m.message[rid] = append(m.message[rid], message...)
}

func (m *MessageInMemory) GetHistory(rid string, msg_id int, offset int) []interface{} {
	var end int
	for i, msg := range m.message[rid] {
		v := reflect.ValueOf(msg)
		if int(reflect.Indirect(v).FieldByName("ID").Int()) == msg_id {
			end = i
			break
		}
	}
	if end-offset < 0 {
		return m.message[rid][:end]
	}
	return m.message[rid][(end - offset - 1):end]
}

func (m *MessageInMemory) GetUnRead(rid string, position int) []interface{} {
	if m.message[rid] != nil && len(m.message[rid]) > position {
		return m.message[rid][position:]
	}
	return nil
}

type TokenInMemory struct {
	token map[string]string
}

func (t *TokenInMemory) Insert(token string, id string) {
	t.token[token] = id
}

func (t *TokenInMemory) Get(token string) string {
	return t.token[token]
}

type IndexInMemory struct {
	roomMember        map[string]map[string]bool
	memberRoom        map[string]map[string]bool
	rommMemberMessage map[string]int
}

func (i *IndexInMemory) SetMemberOfRoom(rid string, mid string) {
	if i.roomMember[rid] == nil {
		i.roomMember[rid] = make(map[string]bool)
	}
	i.roomMember[rid][mid] = true
}

func (i *IndexInMemory) GetMemberFromRoom(rid string) []string {
	res := make([]string, 0)
	for k := range i.roomMember[rid] {
		res = append(res, k)
	}
	return res
}

func (i *IndexInMemory) SetRoomOfMember(mid string, rid string) {
	if i.memberRoom[mid] == nil {
		i.memberRoom[mid] = make(map[string]bool)
	}
	i.memberRoom[mid][rid] = true
}

func (i *IndexInMemory) GetRoomFromMember(mid string) []string {
	res := make([]string, 0)
	for k := range i.memberRoom[mid] {
		res = append(res, k)
	}
	return res
}

func (i *IndexInMemory) SetRoomMemberMessage(rid string, mid string, position int) {
	rmid := rid + ":" + mid
	i.rommMemberMessage[rmid] = position
}

func (i *IndexInMemory) GetRoomMemberMessage(rid string, mid string) int {
	rmid := rid + ":" + mid
	return i.rommMemberMessage[rmid]
}
