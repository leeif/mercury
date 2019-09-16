package redis

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"

	avl "github.com/Workiva/go-datastructures/tree/avl"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-redis/redis"
	"github.com/leeif/mercury/config"
	"github.com/leeif/mercury/storage/data"
	"github.com/leeif/mercury/storage/memory"
)

type Redis struct {
	logger      log.Logger
	client      *redis.Client
	memoryStore *memory.Memory
	msgIDMutex  sync.Mutex
}

func (r *Redis) initRedis(l log.Logger, config config.RedisConfig) {
	var err error
	r.logger = log.With(l, "component", "redis")
	r.client = redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
	})

	pong, err := r.client.Ping().Result()
	r.fatalErr(err)
	level.Debug(r.logger).Log("msg", pong)

	r.memoryStore = &memory.Memory{
		Room:   avl.NewImmutable(),
		Member: avl.NewImmutable(),
	}
}

func (r *Redis) InsertRoomMember(rid string, mid string) {
	res, err := r.client.HMGet("message_pos", rid+"_"+mid).Result()
	r.checkErr(err)
	if len(res) > 0 && res[0] != nil {
		return
	}
	field := make(map[string]interface{})
	field[rid+"_"+mid] = strconv.Itoa(r.getLatestMessage(rid))
	err = r.client.HMSet("message_pos", field).Err()
	r.checkErr(err)
	return
}

func (r *Redis) getLatestMessage(rid string) int {
	sort := &redis.Sort{
		Offset: 0,
		Count:  1,
		Order:  "DESC",
	}
	res, err := r.client.Sort("room_"+string(rid)+"_msg", sort).Result()
	r.checkErr(err)
	if len(res) == 0 {
		return 0
	}
	s, err := strconv.Atoi(res[0])
	r.checkErr(err)
	return s
}

func (r *Redis) InsertMember(member ...interface{}) {
	r.memoryStore.InsertMember(member...)
}

func (r *Redis) GetMember(mid ...string) []interface{} {
	return r.memoryStore.GetMember(mid...)
}

func (r *Redis) InsertToken(mid string, token string) {
	field := make(map[string]interface{})
	field[mid] = token
	err := r.client.HMSet("token", field).Err()
	r.checkErr(err)
}

func (r *Redis) GetToken(mid string) string {
	res, err := r.client.HMGet("token", mid).Result()
	r.checkErr(err)
	return res[0].(string)
}

func (r *Redis) InsertMessage(message *data.MessageBase) int {
	field := make(map[string]interface{})
	r.msgIDMutex.Lock()
	id, err := r.client.Incr("msg_id").Result()
	message.ID = int(id)
	b, err := json.Marshal(message)
	if err != nil {

	}
	field[strconv.Itoa(int(id))] = string(b)
	err = r.client.HMSet("message", field).Err()
	r.checkErr(err)

	err = r.client.LPush("room_"+message.RID+"_msg", strconv.Itoa(int(id))).Err()
	r.checkErr(err)

	r.msgIDMutex.Unlock()
	return int(id)
}

func (r *Redis) GetUnReadMessage(rid string, msgID int) []*data.MessageBase {
	messages := make([]*data.MessageBase, 0)
	sort := &redis.Sort{
		Order: "DESC",
	}
	res, err := r.client.Sort("room_"+rid+"_msg", sort).Result()
	if len(res) == 0 {
		return messages
	}
	r.checkErr(err)
	position := 0
	for i, id := range res {
		if id == strconv.Itoa(msgID) {
			position = i
			break
		}
	}
	if len(res) == 0 {
		return messages
	}
	for _, id := range res[:position] {
		message := &data.MessageBase{}
		res, err := r.client.HMGet("message", id).Result()
		r.checkErr(err)
		if len(res) == 0 || res[0] == nil {
			continue
		}
		err = json.Unmarshal([]byte(res[0].(string)), &message)
		if err != nil {
			continue
		}
		messages = append([]*data.MessageBase{message}, messages...)
	}
	return messages
}

func (r *Redis) GetHistoryMessage(rid string, msgID int, offset int) []*data.MessageBase {
	sort := &redis.Sort{
		Order: "DESC",
	}
	position := -1
	start, end := 0, 0
	res, err := r.client.Sort("room_"+rid+"_msg", sort).Result()
	r.checkErr(err)
	for i, id := range res {
		if id == string(msgID) {
			position = i
			break
		}
	}
	if position-offset < 0 {
		start, end = 0, position
	} else {
		start, end = position-offset, position
	}
	messages := make([]*data.MessageBase, 0)
	for _, id := range res[start:end] {
		message := &data.MessageBase{}
		res, err := r.client.HMGet("message", id).Result()
		r.checkErr(err)
		if len(res) == 0 || res[0] == nil {
			continue
		}
		err = json.Unmarshal([]byte(res[0].(string)), res[0])
		if err != nil {
			continue
		}
		messages = append(messages, message)
	}
	return messages
}

func (r *Redis) SetMemberOfRoom(rid string, mid string) {
	err := r.client.SAdd("room_"+rid+"_member", mid).Err()
	r.checkErr(err)
}

func (r *Redis) GetMemberFromRoom(rid string) []string {
	res, err := r.client.SMembers("room_" + rid + "_member").Result()
	r.checkErr(err)
	return res
}

func (r *Redis) SetRoomOfMember(mid string, rid string) {
	err := r.client.SAdd("member_"+mid+"_room", rid).Err()
	r.checkErr(err)
}

func (r *Redis) GetRoomFromMember(mid string) []string {
	res, err := r.client.SMembers("member_" + mid + "_room").Result()
	r.checkErr(err)
	return res
}

func (r *Redis) SetRoomMemberMessage(rid string, mid string, msgID int) {
	field := make(map[string]interface{})
	field[rid+"_"+mid] = strconv.Itoa(msgID)
	err := r.client.HMSet("message_pos", field).Err()
	r.checkErr(err)
}

func (r *Redis) GetRoomMemberMessage(rid string, mid string) int {
	res, err := r.client.HMGet("message_pos", rid+"_"+mid).Result()
	r.checkErr(err)
	if res[0] == nil {
		return 0
	}
	s, err := strconv.Atoi(res[0].(string))
	r.checkErr(err)
	return s
}

func (r *Redis) checkErr(err error) {
	if err != nil {
		level.Error(r.logger).Log("msg", err)
	}
}

func (r *Redis) fatalErr(err error) {
	if err != nil {
		level.Error(r.logger).Log("msg", err)
		os.Exit(1)
	}
}

func NewRedis(l log.Logger, config config.RedisConfig) *Redis {
	redis := &Redis{}
	redis.initRedis(l, config)
	return redis
}
