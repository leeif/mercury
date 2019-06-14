package server

import (
	"net"
	"net/http"
	"os"
	"regexp"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/leeif/mercury/common"
	h "github.com/leeif/mercury/house"
	"github.com/pkg/errors"
)

var (
	house  *h.House
	logger log.Logger
)

type ServerConfig struct {
	APIAddress *Address
	APIPort    *Port
	WSAddress  *Address
	WSPort     *Port
}

type Address struct {
	s   string
	ip  net.IP
	net *net.IPNet
}

func (a *Address) Set(s string) error {
	var re *regexp.Regexp
	ipString := s
	re, _ = regexp.Compile(`^(\d+\.\d+.\d+.\d+)(|/\d+)$`)
	match := re.FindStringSubmatch(ipString)
	if len(match) == 0 {
		return errors.Errorf("wrong format of address %q", s)
	}
	// not have a range, add /32
	if match[2] == "" {
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
	return a.net.String()
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

func Serve(config *ServerConfig, h *h.House, l log.Logger, exitCh chan error) {
	logger = log.With(l, "component", "server")
	house = h

	waitGroup := &common.WaitGroupWrapper{}

	apiRouter := newAPIRouter()

	apiListener, err := net.Listen("tcp4", ":"+config.APIPort.String())

	if err != nil {
		level.Error(logger).Log("error", err.Error())
		os.Exit(1)
	}

	apiServer := &http.Server{
		Handler: apiRouter,
	}

	waitGroup.Wrap(func() {
		level.Info(logger).Log("msg", "API server is listening at "+config.APIAddress.String()+":"+config.APIPort.String())
		err := apiServer.Serve(apiListener)
		if err != nil {
			exitCh <- err
		}
	})

	wsRouter := newWSRouter()

	wsListener, err := net.Listen("tcp4", ":"+config.WSPort.String())

	if err != nil {
		level.Error(logger).Log("error", err.Error())
		os.Exit(1)
	}

	wsServer := &http.Server{
		Handler: wsRouter,
	}

	waitGroup.Wrap(func() {
		level.Info(logger).Log("msg", "WebSocket server is listening at "+config.WSAddress.String()+":"+config.WSPort.String())
		err := wsServer.Serve(wsListener)
		if err != nil {
			exitCh <- err
		}
	})
}
