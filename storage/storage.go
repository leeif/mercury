package storage

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/leeif/mercury/storage/config"
	"github.com/leeif/mercury/storage/data"
	"github.com/leeif/mercury/storage/memory"
	"github.com/leeif/mercury/storage/mysql"
	"github.com/leeif/mercury/storage/redis"
)

var (
	logger log.Logger
)

func NewStore(l log.Logger, conf *config.StorageConfig) data.Store {
	logger = log.With(l, "component", "storage")
	var store data.Store
	if conf.RedisConfig.Host != "" {
		level.Info(logger).Log("msg", "Using Redis storage")
		store = newRedis(logger, conf.RedisConfig)
	} else if conf.MySQLConfig.Host != "" {
		level.Info(logger).Log("msg", "Using MySQL storage")
		store = newMysql(logger, conf.MySQLConfig)
	} else {
		store = newMemory(logger)
	}
	return store
}

func newMysql(l log.Logger, config *config.MySQLConfig) data.Store {
	store := mysql.NewMySQL(l, config)
	return store
}

func newRedis(l log.Logger, config *config.RedisConfig) data.Store {
	store := redis.NewRedis(l, config)
	return store
}

func newMemory(l log.Logger) data.Store {
	store := memory.NewMemory(l)
	return store
}
