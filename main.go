package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	var port = flag.String("port", "9090", "chat server port")
	var host = flag.String("host", "localhost", "chat server host address")
	addr := *host + ":" + *port

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		route.r(w, r)
	})

	log.Fatal(http.ListenAndServe(addr, nil))
}
