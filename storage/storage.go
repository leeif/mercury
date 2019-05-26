package storage

import (
	avl "github.com/Workiva/go-datastructures/tree/avl"
	"github.com/leeif/mercury/storage/data"
	"github.com/leeif/mercury/storage/memory"
	"github.com/leeif/mercury/storage/config"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func SetLogFlag(a *kingpin.Application, config *config.StorageConfig) {
	a.Flag("mysql.host", "").
		Default("").StringVar(&config.MySQLConfig.Host)

	a.Flag("mysql.port", "").
		Default("").StringVar(&config.MySQLConfig.Port)

	a.Flag("mysql.user", "").
		Default("").StringVar(&config.MySQLConfig.User)

	a.Flag("mysql.password", "").
		Default("").StringVar(&config.MySQLConfig.Password)

	a.Flag("redis.host", "").
		Default("").StringVar(&config.RedisConfig.Host)

	a.Flag("redis.port", "").
		Default("").StringVar(&config.RedisConfig.Port)

	a.Flag("redis.user", "").
		Default("").StringVar(&config.RedisConfig.User)

	a.Flag("redis.password", "").
		Default("").StringVar(&config.RedisConfig.Password)
}

type Store struct {
	Member  data.Member
	Room    data.Room
	Message data.Message
	Token   data.Token
	Index   data.Index
}

func NewStore(config *config.StorageConfig) *Store {
	store := &Store{}
	if config.RedisHost != "" {

	} else if config.MysqlHost != "" {

	} else {
		newMemory(store)
	}
	return store
}

func newMysql(store *Store, config *config.StorageConfig) {

}

func newMemory(store *Store) {
	store.Room = &memory.RoomInMemory{
		Room: avl.NewImmutable(),
	}

	store.Member = &memory.MemberInMemory{
		Member: avl.NewImmutable(),
	}

	store.Message = &memory.MessageInMemory{
		Message: make(map[string][]interface{}),
	}

	store.Token = &memory.TokenInMemory{
		Token: make(map[string]string),
	}

	store.Index = &memory.IndexInMemory{
		MemberRoom:        make(map[string]map[string]bool),
		RoomMember:        make(map[string]map[string]bool),
		RommMemberMessage: make(map[string]int),
	}
}
