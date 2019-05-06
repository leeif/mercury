package house

import (
	"encoding/json"
	"errors"
	"time"
	"mercury/storage"
)

type Message struct {
	storage.MessageBase
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

	// message create timestamp(s)
	msg.CreateTime = time.Now().Unix()
	if msgMap["msgType"] != nil {
		msg.MsgType = int(msgMap["msgType"].(float64))
	} else {
		return nil, errors.New("No msgType")
	}

	if msgMap["mid"] != nil {
		msg.MID = msgMap["mid"].(string)
	}

	if msgMap["rid"] != nil {
		msg.RID = msgMap["rid"].(string)
	}

	if msgMap["text"] != nil {
		msg.Text = msgMap["text"].(string)
	}

	return msg, nil
}
