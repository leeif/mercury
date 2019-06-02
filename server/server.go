package server

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	h "github.com/leeif/mercury/house"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	house  *h.House
	route  *Route
	logger log.Logger
)

type ServerConfig struct {
	Address *Address
	Port    *Port
}

type Address struct {
	s string
}

func (a *Address) Set(s string) error {
	a.s = s
	return nil
}

func (a *Address) String() string {
	return a.s
}

type Port struct {
	s string
}

func (p *Port) Set(s string) error {
	p.s = s
	return nil
}

func (p *Port) String() string {
	return p.s
}

func SetServerFlag(a *kingpin.Application, config *ServerConfig) {
	config.Address = &Address{}
	a.Flag("server.address", "server listen address").
		Default("127.0.0.1").SetValue(config.Address)

	config.Port = &Port{}
	a.Flag("server.port", "server listen port").
		Default("6010").SetValue(config.Port)
}

func Serve(config *ServerConfig, h *h.House, l log.Logger) {
	logger = log.With(l, "component", "server")
	house = h
	rt := NewRoute()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rt.Select(w, r)
	})

	addr := config.Address.s + ":" + config.Port.s
	level.Debug(logger).Log("listen", addr)
	level.Error(logger).Log(http.ListenAndServe(addr, nil))
}
