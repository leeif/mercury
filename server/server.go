package server

import (
	"net"
	"net/http"
	"regexp"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	h "github.com/leeif/mercury/house"
	"github.com/pkg/errors"
)

var (
	house  *h.House
	route  *Route
	logger log.Logger
)

type ServerConfig struct {
	APIAddress *Address
	WSAddress  *Address
	Port       *Port
}

type Address struct {
	s   string
	ip  net.IP
	net *net.IPNet
}

func (a *Address) Set(s string) error {
	var re *regexp.Regexp
	ipString := s
	re, _ = regexp.Compile(`^([0-9]+\.){3}([0-9])/([0-9]+)$`)
	maskFormat := re.MatchString(ipString)
	re, _ = regexp.Compile(`^([0-9]+\.){3}([0-9])$`)
	ipFormat := re.MatchString(ipString)
	if !maskFormat && !ipFormat {
		return errors.Errorf("wrong format of address %q", s)
	} else if !maskFormat && ipFormat {
		ipString += "/32"
	}
	ip, net, err := net.ParseCIDR(ipString)
	if err != nil {
		return errors.Errorf("wrong format of address %q", s)
	}
	a.s = ipString
	a.ip = ip
	a.net = net
	return nil
}

func (a *Address) String() string {
	return a.s
}

func (a *Address) Contains(s string) bool {
	ip := net.ParseIP(s)
	return a.net.Contains(ip)
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

func Serve(config *ServerConfig, h *h.House, l log.Logger) {
	logger = log.With(l, "component", "server")
	house = h
	rt := NewRoute(config)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rt.Select(w, r)
	})

	addr := "0.0.0.0:" + config.Port.s
	level.Info(logger).Log("msg", "Rest api address: " + config.APIAddress.s)
	level.Info(logger).Log("msg", "Websocket address: " + config.WSAddress.s)
	level.Info(logger).Log("msg", "Port: " + config.Port.s)
	level.Error(logger).Log(http.ListenAndServe(addr, nil))
}
