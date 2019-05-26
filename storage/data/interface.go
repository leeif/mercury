package data

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
