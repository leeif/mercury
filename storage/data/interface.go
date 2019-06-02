package data

type Store interface {
	// member
	InsertMember(...interface{})
	GetMember(...string) []interface{}

	// room
	InsertRoom(...interface{})
	InsertRoomMember(interface{}, interface{})
	GetRoom(...string) []interface{}

	// token
	InsertToken(string, string)
	GetToken(string) string

	// message
	InsertMessage(*MessageBase) int
	GetUnReadMessage(string, int) []*MessageBase
	GetHistoryMessage(string, int, int) []*MessageBase

	// index
	SetMemberOfRoom(string, string)
	GetMemberFromRoom(string) []string
	SetRoomOfMember(string, string)
	GetRoomFromMember(string) []string
	SetRoomMemberMessage(string, string, int)
	GetRoomMemberMessage(string, string) int
}
