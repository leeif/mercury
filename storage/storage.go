package storage

import (
	avl "github.com/Workiva/go-datastructures/tree/avl"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/leeif/mercury/storage/config"
	"github.com/leeif/mercury/storage/data"
	"github.com/leeif/mercury/storage/memory"
	"github.com/leeif/mercury/storage/mysql"
)

var (
	logger log.Logger
)

func NewStore(l log.Logger, conf *config.StorageConfig) data.Store {
	logger = log.With(l, "component", "storage")
	var store data.Store
	if conf.RedisConfig.Host != "" {

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

func newMemory(l log.Logger) data.Store {
	store := &memory.Memory{
		Room:              avl.NewImmutable(),
		Member:            avl.NewImmutable(),
		Message:           make(map[string][]*data.MessageBase),
		Token:             make(map[string]string),
		MemberRoom:        make(map[string]map[string]bool),
		RoomMember:        make(map[string]map[string]bool),
		RommMemberMessage: make(map[string]int),
	}
	return store
}
