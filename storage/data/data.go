package data

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
	ID int `json:"id"`
	// from member id
	MID string `json:"mid"`
	// room id
	RID        string `json:"rid"`
	CreateTime int64  `json:"createTime"`
	Text       string `json:"text"`
}

type RoomMemberIndex struct {
	RID string
	MID string
}

func (rmi *RoomMemberIndex) Compare(entry avl.Entry) int {
	other := entry.(*RoomMemberIndex)
	if rmi.RID == other.RID {
		if other.MID == "" || rmi.MID == other.MID {
			return 0
		} else if rmi.MID > other.MID {
			return 1
		} else {
			return -1
		}
	} else if rmi.RID > other.RID {
		return 1
	} else {
		return -1
	}
}

type MemberRoomIndex struct {
	MID string
	RID string
}

func (mri *MemberRoomIndex) Compare(entry avl.Entry) int {
	other := entry.(*MemberRoomIndex)
	if mri.MID == other.MID {
		if other.RID == "" || mri.RID == other.RID {
			return 0
		} else if mri.RID > other.RID {
			return 1
		} else {
			return -1
		}
	} else if mri.MID > other.MID {
		return 1
	} else {
		return -1
	}
}

type RoomMemberMessageIndex struct {
	// id of "roomid-memberid"
	RMID  string
	Msgid string
}

func (rmmi *RoomMemberMessageIndex) Compare(entry avl.Entry) int {
	other := entry.(*RoomMemberMessageIndex)
	if rmmi.RMID == other.RMID {
		if other.Msgid == "" || rmmi.Msgid == other.Msgid {
			return 0
		} else if rmmi.Msgid > other.Msgid {
			return 1
		} else {
			return -1
		}
	} else if rmmi.RMID > other.RMID {
		return 1
	} else {
		return -1
	}
}
