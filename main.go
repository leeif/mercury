package main

import (
	"fmt"
	"os"
	"path/filepath"
	_ "github.com/go-sql-driver/mysql"
	conf "github.com/leeif/mercury/config"
	"github.com/leeif/mercury/common"
	"github.com/leeif/mercury/storage"
	c "github.com/leeif/mercury/connection"
	h "github.com/leeif/mercury/house"
	"github.com/leeif/mercury/server"
	"github.com/pkg/errors"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	config := conf.Config{
		LogConfig:    common.LogConfig{},
		ServerConfig: server.ServerConfig{},
	}

	a := kingpin.New(filepath.Base(os.Args[0]), "Mercury server")
	a.HelpFlag.Short('h')

	// config file path
	a.Flag("config.file", "configure file path").Default("mc.conf").StringVar(&config.ConfigFile)
	
	// load server around command line option
	server.SetServerFlag(a, &config.ServerConfig)

	// load log around command line option
	common.SetLogFlag(a, &config.LogConfig)

	// load storage around command line option
	storage.SetLogFlag(a, &config.StorageConfig)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing commandline arguments"))
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	logger   := common.NewLogger(&config.LogConfig)
	connPool := c.NewPool(nil, logger)

	s        := storage.NewStore(logger, &config.StorageConfig)
	house    := h.NewHouse(logger, s, connPool)

	server.Serve(&config.ServerConfig, house, logger)
}
