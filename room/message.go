package room

import (
	"encoding/json"
)

// Message type
const (
	// REGISTERROOM : register user into a chat room
	REGISTERROOM = iota
	// SENDTOROOM : send a message to a room
	SENDMESSAGE
)

type Message struct {
	MsgType    int      `json:"-"`
	FromID     string   `json:"fromID"`
	ToID       string   `json:"toID"`
	RoomID     []string `json:"roomID"`
	CreateTime string   `json:"createTime"`
	Body       string   `json:"body"`
}

type Response struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

func (message *Message) json() ([]byte, error) {
	b, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func NewMessage(data []byte) (*Message, error) {
	msgMap := make(map[string]interface{})
	msg := &Message{}
	err := json.Unmarshal(data, &msgMap)
	if err != nil {
		return nil, err
	}

	if msgMap["msgType"] != nil {
		msg.MsgType = int(msgMap["msgType"].(float64))
	}

	if msgMap["toID"] != nil {
		msg.ToID = msgMap["toID"].(string)
	}

	if msgMap["roomID"] != nil {
		roomID := make([]string, len(msgMap["roomID"].([]interface{})))
		for i, v := range msgMap["roomID"].([]interface{}) {
			roomID[i] = v.(string)
		}
		msg.RoomID = roomID
	}

	if msgMap["body"] != nil {
		msg.Body = msgMap["body"].(string)
	}

	return msg, nil
}
