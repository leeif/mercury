package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
	"github.com/gorilla/websocket"
)

const TestRoom = "2333"

var port = flag.String("port", "6010", "chat server port")
var host = flag.String("host", "localhost", "chat server host address")
var path = flag.String("path", "/ws/connect", "ws connection path")
var member = flag.String("member", "1", "member")

var logger log.Logger
var token = ""

func init () {
	logger = log.NewLogfmtLogger(os.Stdin)
	logger = log.With(logger, "caller", log.DefaultCaller, "component", "client")
}

type Message struct {
	Type int    `json:"type"`
	RID  string `json:"rid"`
	MID  string `json:"mid"`
	Text string `json:"text"`
}

func AddRoom() {
	url := "http://" + *host + ":" + *port + "/api/room/add?" + "room=" + TestRoom + "&member=" + *member
	resp, err := http.Post(url, "", nil)
	if err != nil {
		level.Error(logger).Log("Failed")
	} else {
		level.Info(logger).Log("RoomID", TestRoom)
	}
	defer resp.Body.Close()
}

func GetToken() string {
	level.Info(logger).Log("Token", *member)
	url := "http://" + *host + ":" + *port + "/api/token?" + "id=" + *member
	resp, err1 := http.Get(url)
	if err1 != nil {
		level.Error(logger).Log("Get Token Error")
	}
	defer resp.Body.Close()
	b, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		level.Error(logger).Log("Read Resp Body Error")
	}
	res := make(map[string]interface{})
	err3 := json.Unmarshal(b, &res)
	if err3 != nil {
		level.Error(logger).Log("Parse Resp Body Error")
	}
	body := res["body"].(map[string]interface{})
	return body["token"].(string)
}

func main() {
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// api operation
	AddRoom()
	token := GetToken()

	addr := *host + ":" + *port

	u := url.URL{Scheme: "ws", Host: addr, Path: *path}
	u.RawQuery = "token=" + token
	level.Info(logger).Log("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		level.Error(logger).Log("Can not connect " + addr + *path)
	}

	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		fmt.Print("send> ")
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				level.Info(logger).Log("read: %v", err)
				return
			}
			m := Message{}
			err2 := json.Unmarshal(message, &m)
			if err2 == nil {
				if m.MID == *member {
					fmt.Printf("receive> message From member %s : %s", m.MID, m.Text)
				} else {
					fmt.Printf("\nreceive> message From member %s : %s", m.MID, m.Text)
				}
			}
			fmt.Print("send> ")
		}
	}()

	go func() {
		defer close(done)
		reader := bufio.NewReader(os.Stdin)
		for {
			text, _ := reader.ReadString('\n')
			message := Message{Type: 1, RID: TestRoom, Text: text}
			var b []byte
			var err error
			b, err = json.Marshal(message)
			if err != nil {
				continue
			}
			err = c.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				level.Error(logger).Log("write:", err)
				return
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			level.Info(logger).Log("interrupt")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

			if err != nil {
				level.Info(logger).Log("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
