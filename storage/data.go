package storage

import (
	"reflect"

	"github.com/Workiva/go-datastructures/tree/avl"
)

type RoomBase struct {
	ID string
}

func (rb *RoomBase) Compare(entry avl.Entry) int {
	v := reflect.ValueOf(entry)
	id := reflect.Indirect(v).FieldByName("ID").String()
	if rb.ID == id {
		return 0
	} else if rb.ID > id {
		return 1
	} else {
		return -1
	}
}

type MemberBase struct {
	ID string
}

func (mb *MemberBase) Compare(entry avl.Entry) int {
	v := reflect.ValueOf(entry)
	id := reflect.Indirect(v).FieldByName("ID").String()
	if mb.ID == id {
		return 0
	} else if mb.ID > id {
		return 1
	} else {
		return -1
	}
}

type MessageBase struct {
	ID      int
	MsgType int
	// mid[0] to, mid[1] from
	MID []string
	// room id
	RID        []string
	CreateTime int64
	Text       string
}

type RoomMemberIndex struct {
	rid string
	mid string
}

func (rmi *RoomMemberIndex) Compare(entry avl.Entry) int {
	other := entry.(*RoomMemberIndex)
	if rmi.rid == other.rid {
		if other.mid == "" || rmi.mid == other.mid {
			return 0
		} else if rmi.mid > other.mid {
			return 1
		} else {
			return -1
		}
	} else if rmi.rid > other.rid {
		return 1
	} else {
		return -1
	}
}

type MemberRoomIndex struct {
	mid string
	rid string
}

func (mri *MemberRoomIndex) Compare(entry avl.Entry) int {
	other := entry.(*MemberRoomIndex)
	if mri.mid == other.mid {
		if other.rid == "" || mri.rid == other.rid {
			return 0
		} else if mri.rid > other.rid {
			return 1
		} else {
			return -1
		}
	} else if mri.mid > other.mid {
		return 1
	} else {
		return -1
	}
}

type RoomMemberMessageIndex struct {
	// id of "roomID-memberID"
	rmid  string
	msgid string
}

func (rmmi *RoomMemberMessageIndex) Compare(entry avl.Entry) int {
	other := entry.(*RoomMemberMessageIndex)
	if rmmi.rmid == other.rmid {
		if other.msgid == "" || rmmi.msgid == other.msgid {
			return 0
		} else if rmmi.msgid > other.msgid {
			return 1
		} else {
			return -1
		}
	} else if rmmi.rmid > other.rmid {
		return 1
	} else {
		return -1
	}
}
