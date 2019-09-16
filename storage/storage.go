package storage

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/leeif/mercury/config"
	"github.com/leeif/mercury/storage/data"
	"github.com/leeif/mercury/storage/memory"
	"github.com/leeif/mercury/storage/mysql"
	"github.com/leeif/mercury/storage/redis"
)

var (
	logger log.Logger
)

func NewStore(l log.Logger, cfg *config.Config) data.Store {
	scfg := cfg.Storage
	logger = log.With(l, "component", "storage")
	var store data.Store
	switch scfg.Type.String() {
	case "memory":
		level.Info(logger).Log("msg", "Using Memory storage")
		store = newMemory(logger)
	case "mysql":
		level.Info(logger).Log("msg", "Using MySQL storage")
		store = newMysql(logger, scfg.MySQLConfig)
	case "redis":
		level.Info(logger).Log("msg", "Using Redis storage")
		store = newRedis(logger, scfg.RedisConfig)
	}
	return store
}

func newMysql(l log.Logger, config config.MySQLConfig) data.Store {
	store := mysql.NewMySQL(l, config)
	return store
}

func newRedis(l log.Logger, config config.RedisConfig) data.Store {
	store := redis.NewRedis(l, config)
	return store
}

func newMemory(l log.Logger) data.Store {
	store := memory.NewMemory(l)
	return store
}
