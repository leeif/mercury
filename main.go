package main

import (
	"flag"
	"log"
	"net/http"
	room "mercury/room"
)

func main() {
	var port = flag.String("port", "9090", "chat server port")
	var host = flag.String("host", "localhost", "chat server host address")
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		member := room.NewMember()
		member.GenerateConnection(w, r)
	})

	addr := *host + ":" + *port
	log.Fatal(http.ListenAndServe(addr, nil))
}
