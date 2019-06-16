package config

import (
	log "github.com/leeif/mercury/common"
	"github.com/leeif/mercury/server"
	storage "github.com/leeif/mercury/storage/config"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func AddFlag(a *kingpin.Application, c *Config) {
	// config file path
	a.Flag("config.file", "configure file path").Default("mc.cnf.toml").StringVar(&c.ConfigFile)

	SetLogFlag(a, &c.Log)
	SetServerFlag(a, &c.Server)
	SetStorageFlag(a, &c.Storage)
}

func SetServerFlag(a *kingpin.Application, config *server.ServerConfig) {
	config.APIAddress = &server.Address{}
	a.Flag("server.api.address", "server listen address").SetValue(config.APIAddress)

	config.WSAddress = &server.Address{}
	a.Flag("server.ws.address", "server listen address").SetValue(config.WSAddress)

	config.APIPort = &server.Port{}
	a.Flag("server.api.port", "server listen port").SetValue(config.APIPort)

	config.WSPort = &server.Port{}
	a.Flag("server.ws.port", "server listen port").SetValue(config.WSPort)	
}

func SetStorageFlag(a *kingpin.Application, conf *storage.StorageConfig) {
	conf.MySQLConfig = &storage.MySQLConfig{}
	conf.RedisConfig = &storage.RedisConfig{}

	a.Flag("storage.maxunread", "max unread which will be send when client connected").Default("-1").IntVar(&conf.MaxUnRead)

	a.Flag("mysql.host", "mysql host").StringVar(&conf.MySQLConfig.Host)

	a.Flag("mysql.port", "mysql port").StringVar(&conf.MySQLConfig.Port)

	a.Flag("mysql.user", "mysql user").StringVar(&conf.MySQLConfig.User)

	a.Flag("mysql.password", "mysql possword").StringVar(&conf.MySQLConfig.Password)

	a.Flag("redis.host", "redis host").StringVar(&conf.RedisConfig.Host)

	a.Flag("redis.port", "redis port").StringVar(&conf.RedisConfig.Port)

	a.Flag("redis.password", "redis password").StringVar(&conf.RedisConfig.Password)

	a.Flag("redis.db", "redis db number").Default("-1").IntVar(&conf.RedisConfig.DB)
}

func SetLogFlag(a *kingpin.Application, config *log.LogConfig) {
	config.Level = &log.AllowedLevel{}
	a.Flag("log.level", "[debug, info, warn, error]").SetValue(config.Level)

	config.Format = &log.AllowedFormat{}
	a.Flag("log.format", "[logfmt, json]").SetValue(config.Format)
}
