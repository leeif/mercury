package storage

import (
	avl "github.com/Workiva/go-datastructures/tree/avl"
)

type Member interface {
	Insert(...interface{})
	Get(...string) []avl.Entry
}

type Room interface {
	Insert(...interface{})
	Get(...string) []avl.Entry
}

type Token interface {
	Insert(string, string)
	Get(string) string
}

type Message interface {
	Insert(string, ...interface{})
	GetUnRead(string, int) []interface{}
	GetHistory(string, int, int) []interface{}
}

type Index interface {
	SetMemberOfRoom(string, string)
	GetMemberFromRoom(string) []string
	SetRoomOfMember(string, string)
	GetRoomFromMember(string) []string
	SetRoomMemberMessage(string, string, int)
	GetRoomMemberMessage(string, string) int
}

type Store struct {
	Member  Member
	Room    Room
	Message Message
	Token   Token
	Index   Index
}

func NewStore() *Store {
	return &Store {
		Room:    newRoom(),
		Member:  newMember(),
		Message: newMessage(),
		Token:   newToken(),
		Index:   newIndex(),
	}
}

func newRoom() Member {
	return &RoomInMemory {
		room: avl.NewImmutable(),
	}
}

func newMember() Member {
	return &MemberInMemory {
		member: avl.NewImmutable(),
	}
}

func newMessage() Message {
	return &MessageInMemory {
		message: make(map[string][]interface{}),
	}
}

func newIndex() Index {
	return &IndexInMemory {
		memberRoom:        make(map[string]map[string]bool),
		roomMember:        make(map[string]map[string]bool),
		rommMemberMessage: make(map[string]int),
	}
}

func newToken() Token {
	return &TokenInMemory {
		token: make(map[string]string),
	}
}