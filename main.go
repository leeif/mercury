package main

import (
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/leeif/mercury/common"
	conf "github.com/leeif/mercury/config"
	conn "github.com/leeif/mercury/connection"
	house "github.com/leeif/mercury/house"
	"github.com/leeif/mercury/server"
	"github.com/leeif/mercury/storage"
	"github.com/pkg/errors"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	config := conf.Config{}

	a := kingpin.New(filepath.Base(os.Args[0]), "Mercury server")
	a.HelpFlag.Short('h')

	// config file path
	a.Flag("config.file", "configure file path").Default("mc.cnf.toml").StringVar(&config.ConfigFile)

	// load flag
	conf.AddFlag(a, &config)
	_, err := a.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing commandline arguments"))
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	// load configure file
	conf.LoadConfigFile(config.ConfigFile, &config)

	logger := common.NewLogger(&config.Log)
	connPool := conn.NewPool(nil, logger)

	store := storage.NewStore(logger, &config.Storage)
	house := house.NewHouse(logger, store, connPool)

	server.Serve(&config.Server, house, logger)
}
