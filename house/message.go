package house

import (
	"encoding/json"
	"errors"
	"time"
	"github.com/leeif/mercury/storage/data"
)

const (
	SEND          = 1
	HISTORY       = 2
	DEFAULTOFFEST = 10
)

var (
	msgID = 0
)

type History struct {
	// room id
	RID      string
	// start msg id
	MsgID    int
	Offest   int
}

type Message struct {
	data.MessageBase
}

type Response struct {
	Status string         `json:"status"`
	Body   interface{}    `json:"body"`
}

func (reponse *Response) json() ([]byte, error) {
	b, err := json.Marshal(reponse)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (message *Message) json() ([]byte, error) {
	b, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func newMessage(data []byte) (int, interface{}, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(data, &m)
	if err != nil {
		return -1, nil, err
	}
	var t int
	if m["type"] != nil {
		t = int(m["type"].(float64))
	} else {
		return -1, nil, errors.New("No Type")
	}
	switch t {
	case SEND:
		item, err := newSend(m)
		return t, item, err
	case HISTORY:
		item, err := newHisotry(m)
		return t, item, err
	}
	return -1, nil, errors.New("Empty")
}

func newHisotry(historyMap map[string]interface{}) (*History, error) {
	history := &History{}

	if historyMap["rid"] != nil {
		history.RID = historyMap["rid"].(string)
	} else {
		return nil, errors.New("No Room ID")
	}

	if historyMap["msg_id"] != nil {
		history.MsgID = int(historyMap["msg_id"].(float64))
	} else {
		return nil, errors.New("No Start Msg ID")
	}

	if historyMap["offest"] != nil {
		history.Offest = int(historyMap["offest"].(float64))
	} else {
		// default
		history.Offest = DEFAULTOFFEST
	}

	return history, nil
}

func newSend(msgMap map[string]interface{}) (*Message, error) {
	msg := &Message{}
	// message create timestamp(s)
	msg.CreateTime = time.Now().Unix()

	if msgMap["mid"] != nil {
		msg.MID = msgMap["mid"].(string)
	}

	if msgMap["rid"] != nil {
		msg.RID = msgMap["rid"].(string)
	}

	msgID++
	msg.ID = msgID

	if msgMap["text"] != nil {
		msg.Text = msgMap["text"].(string)
	}

	return msg, nil
}
