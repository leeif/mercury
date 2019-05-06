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
	Insert(...interface{})
	Get(...string) []avl.Entry
}

type Index interface {
	SetMemberOfRoom(string, string)
	GetMemberFromRoom(string) []string
	SetRoomOfMember(string, string)
	GetRoomFromMember(string) []string
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
		// Message: newMessage(),
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

// func newMessage() Message {
// 	return &MessageInMemory {
// 		message: avl.NewImmutable(),
// 	}
// }

func newIndex() Index {
	return &IndexInMemory {
		memberRoom:        make(map[string]map[string]bool),
		roomMember:        make(map[string]map[string]bool),
		rommMemberMessage: make(map[string]string),
	}
}

func newToken() Token {
	return &TokenInMemory {
		token: make(map[string]string),
	}
}