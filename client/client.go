package main

import (
	"bufio"
	"flag"
	"net/url"
	"os"
	"os/signal"
	"time"
	"mercury/utils"
	"github.com/gorilla/websocket"
)

var port = flag.String("port", "9090", "chat server port")
var host = flag.String("host", "localhost", "chat server host address")
var path = flag.String("path", "/ws/connect", "ws connection path")
var token = flag.String("token", "", "chat token")

func main() {
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	addr := *host + ":" + *port

	u := url.URL{Scheme: "ws", Host: addr, Path: *path}
	u.RawQuery = "token=" + *token
	utils.Info("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		utils.Error("Can not connect " + addr + *path)
	}

	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				utils.Info("read:", err)
				return
			}
			utils.Info("recv: %s", message)
		}
	}()

	go func() {
		defer close(done)
		reader := bufio.NewReader(os.Stdin)
		for {
			text, _ := reader.ReadString('\n')
			utils.Info("send: " + text)
			err := c.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				utils.Info("write:", err)
				return
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			utils.Info("interrupt")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

			if err != nil {
				utils.Info("write close:", err)
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
