package main

import (
	"flag"
	"net/http"
	"github.com/leeif/mercury/route"
	"github.com/leeif/mercury/house"
	"github.com/leeif/mercury/connection"
	"github.com/leeif/mercury/common"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"path/filepath"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"fmt"
	"os"
)

func main() {
	var port = flag.String("port", "9090", "chat server port")
	var host = flag.String("host", "localhost", "chat server host address")
	addr := *host + ":" + *port

	a := kingpin.New(filepath.Base(os.Args[0]), "Mercury server")
	logConfig := common.Config{}
	common.SetFlag(a, &logConfig)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing commandline arguments"))
		a.Usage(os.Args[1:])
		os.Exit(2)
	}


	logger := common.NewLogger(&logConfig)

	house := house.NewHouse(logger)
	connection.WithLogger(logger)
	rt := route.New(logger, house)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rt.Select(w, r)
	})

	level.Error(logger).Log(http.ListenAndServe(addr, nil))
}
