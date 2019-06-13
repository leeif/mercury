package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/websocket"
	"github.com/marcusolsson/tui-go"
)

var myDial = &net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
	DualStack: true,
}

// 自定义DialContext
var myDialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
	network = "tcp4" //仅使用ipv4
	//network = "tcp6" //仅使用ipv6
	return myDial.DialContext(ctx, network, addr)
}

var client = &http.Client{
	Transport: &http.Transport{
		DialContext: myDialContext,
	},
}

const TestRoom = "2333"

var port = flag.String("port", "6010", "chat server port")
var host = flag.String("host", "localhost", "chat server host address")
var path = flag.String("path", "/ws/connect", "ws connection path")
var member = flag.String("member", "1", "member")

var logger log.Logger
var token = ""

func init() {
	logger = log.NewLogfmtLogger(os.Stdin)
	logger = log.With(logger, "caller", log.DefaultCaller, "component", "client")
}

type Message struct {
	Type   int    `json:"type"`
	RID    string `json:"rid"`
	MID    string `json:"mid"`
	Text   string `json:"text"`
	Offset int    `json:"offset"`
	MsgID  int    `json:"msgid"`
}

func AddRoom() {
	url := "http://" + *host + ":" + *port + "/api/room/add?" + "room=" + TestRoom + "&member=" + *member
	resp, err := client.Post(url, "", nil)
	if err != nil {
		level.Error(logger).Log("error", err.Error())
		os.Exit(1)
	} else {
		level.Info(logger).Log("RoomID", TestRoom)
	}
	defer resp.Body.Close()
}

func GetToken() string {
	level.Info(logger).Log("Member", *member)
	url := "http://" + *host + ":" + *port + "/api/token?" + "id=" + *member
	resp, err1 := client.Get(url)
	if err1 != nil {
		level.Error(logger).Log("error", err1.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()
	b, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		level.Error(logger).Log("error", "Read Resp Body Error")
	}
	res := make(map[string]interface{})
	err3 := json.Unmarshal(b, &res)
	if err3 != nil {
		level.Error(logger).Log("error", "Parse Resp Body Error")
	}
	body := res["body"].(map[string]interface{})
	return body["token"].(string)
}

type Connection struct {
	Send chan Message
	Recv chan Message
	ws   *websocket.Conn
}

func (c *Connection) Init() {
	c.Send = make(chan Message, 2)
	c.Recv = make(chan Message, 2)
	AddRoom()
	token := GetToken()

	addr := *host + ":" + *port

	u := url.URL{Scheme: "ws", Host: addr, Path: *path}
	u.RawQuery = "token=" + token
	level.Info(logger).Log("msg", "connecting to "+u.String())

	d := websocket.Dialer{
		NetDialContext: myDialContext,
	}
	var err error
	c.ws, _, err = d.Dial(u.String(), nil)
	if err != nil {
		level.Error(logger).Log("Can not connect " + addr + *path)
		os.Exit(1)
	}
	go c.StartSend()
	go c.StartRecv()
}

func (c *Connection) StartSend() {
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// send failed
			}
			var b []byte
			var err error
			b, err = json.Marshal(message)
			if err != nil {
				continue
			}
			err = c.ws.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				level.Error(logger).Log("write:", err)
			}
		}
	}
}

func (c *Connection) StartRecv() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			level.Info(logger).Log("read: %v", err)
			return
		}
		m := Message{}
		err2 := json.Unmarshal(message, &m)
		if err2 == nil {
			c.Recv <- m
		}
	}
}

func (c *Connection) Close() {
	err := c.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		level.Info(logger).Log("write close:", err)
		return
	}
	c.ws.Close()
}

func appendHistory(message Message, history *tui.Box) {
	now := time.Now()
	history.Append(tui.NewHBox(
		tui.NewLabel(now.Format("15:04")),
		tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<user%s>", message.MID))),
		tui.NewLabel(message.Text),
		tui.NewSpacer(),
	))
}

func main() {
	flag.Parse()

	c := &Connection{}

	sidebar := tui.NewVBox(
		tui.NewLabel("CHANNELS"),
		tui.NewLabel("#"+TestRoom),
		tui.NewSpacer(),
	)
	sidebar.SetBorder(true)

	history := tui.NewVBox()

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	input.OnSubmit(func(e *tui.Entry) {
		message := Message{}
		message.Text = e.Text()
		message.RID = TestRoom
		message.MID = *member
		message.Type = 1
		// send message through websocket
		c.Send <- message

		appendHistory(message, history)
		input.SetText("")
	})

	root := tui.NewHBox(sidebar, chat)

	ui, err := tui.New(root)
	if err != nil {
		level.Debug(logger).Log("error", err)
	}

	ui.SetKeybinding("Esc", func() {
		ui.Quit()

	})

	c.Init()
	defer c.ws.Close()

	go func() {
		for {
			select {
			case message, ok := <-c.Recv:
				if !ok {
					// send failed
					level.Debug(logger).Log("notok")
					continue
				}
				if message.MID != *member {
					ui.Update(func() {
						appendHistory(message, history)
					})
				}
			}
		}
	}()

	if err := ui.Run(); err != nil {
		level.Debug(logger).Log("error", err)
	}
}
