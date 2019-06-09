package config

import (
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"github.com/leeif/mercury/server"
	log "github.com/leeif/mercury/common"
	storage "github.com/leeif/mercury/storage/config"
)

func AddFlag(a *kingpin.Application, c *Config) {
	SetLogFlag(a, &c.Log)
	SetServerFlag(a, &c.Server)
	SetStorageFlag(a, &c.Storage)
}

func SetServerFlag(a *kingpin.Application, config *server.ServerConfig) {
	config.APIAddress = &server.Address{}
	a.Flag("server.api.address", "server listen address").
		Default(DefaultAPIAddress).SetValue(config.APIAddress)

	config.WSAddress = &server.Address{}
	a.Flag("server.ws.address", "server listen address").
		Default(DefaultWSAddress).SetValue(config.WSAddress)

	config.Port = &server.Port{}
	a.Flag("server.port", "server listen port").
		Default(DefaultServerPort).SetValue(config.Port)
}

func SetStorageFlag(a *kingpin.Application, conf *storage.StorageConfig) {
	conf.MySQLConfig = &storage.MySQLConfig{}
	conf.RedisConfig = &storage.RedisConfig{}
	a.Flag("mysql.host", "").
		Default("").StringVar(&conf.MySQLConfig.Host)

	a.Flag("mysql.port", "").
		Default(DefaultMySQLPort).StringVar(&conf.MySQLConfig.Port)

	a.Flag("mysql.user", "").
		Default(DefaultMySQLUser).StringVar(&conf.MySQLConfig.User)

	a.Flag("mysql.password", "").
		Default("").StringVar(&conf.MySQLConfig.Password)

	a.Flag("redis.host", "").
		Default("").StringVar(&conf.RedisConfig.Host)

	a.Flag("redis.port", "").
		Default("").StringVar(&conf.RedisConfig.Port)

	a.Flag("redis.user", "").
		Default("").StringVar(&conf.RedisConfig.User)

	a.Flag("redis.password", "").
		Default("").StringVar(&conf.RedisConfig.Password)
}

func SetLogFlag(a *kingpin.Application, config *log.LogConfig) {
	config.Level = &log.AllowedLevel{}
	a.Flag("log.level", "[debug, info, warn, error]").
		Default(DefaultLogLevel).SetValue(config.Level)

	config.Format = &log.AllowedFormat{}
	a.Flag("log.format", "[logfmt, json]").
		Default(DefaultLogFormat).SetValue(config.Format)
}
