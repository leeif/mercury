package storage

import (
	"github.com/go-kit/kit/log/level"
	avl "github.com/Workiva/go-datastructures/tree/avl"
	"github.com/go-kit/kit/log"
	"github.com/leeif/mercury/storage/config"
	"github.com/leeif/mercury/storage/data"
	"github.com/leeif/mercury/storage/memory"
	"github.com/leeif/mercury/storage/mysql"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	logger log.Logger
)

func SetLogFlag(a *kingpin.Application, conf *config.StorageConfig) {
	conf.MySQLConfig = &config.MySQLConfig{}
	conf.RedisConfig = &config.RedisConfig{}
	a.Flag("mysql.host", "").
		Default("").StringVar(&conf.MySQLConfig.Host)

	a.Flag("mysql.port", "").
		Default("").StringVar(&conf.MySQLConfig.Port)

	a.Flag("mysql.user", "").
		Default("").StringVar(&conf.MySQLConfig.User)

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

func NewStore(l log.Logger, conf *config.StorageConfig) data.Store {
	logger = log.With(l, "component", "storage")
	var store data.Store
	if conf.RedisConfig.Host != "" {

	} else if conf.MySQLConfig.Host != "" {
		level.Debug(logger).Log("msg", "MySQL storage")
		store = newMysql(logger, conf.MySQLConfig)
	} else {
		level.Debug(logger).Log("msg", "Memory storage")
		store = newMemory(logger, conf.RedisConfig)
	}
	return store
}

func newMysql(l log.Logger, config *config.MySQLConfig) data.Store {
	store := mysql.NewMySQL(l, config)
	return store
}

func newMemory(l log.Logger, config *config.RedisConfig) data.Store {
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
