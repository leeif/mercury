package redis

import (
	"os"
	"reflect"
	"strconv"
	"sync"

	avl "github.com/Workiva/go-datastructures/tree/avl"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-redis/redis"
	"github.com/leeif/mercury/storage/config"
	"github.com/leeif/mercury/storage/data"
	"github.com/leeif/mercury/storage/memory"
)

type Redis struct {
	logger      log.Logger
	client      *redis.Client
	memoryStore *memory.Memory
	msgIDMutex  sync.Mutex
}

func (r *Redis) initRedis(l log.Logger, config *config.RedisConfig) {
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

func (r *Redis) InsertRoomMember(room interface{}, member interface{}) {
	var v reflect.Value
	v = reflect.ValueOf(room)
	roomBase := reflect.Indirect(v).FieldByName("RoomBase").Interface().(data.RoomBase)
	v = reflect.ValueOf(member)
	memberBase := reflect.Indirect(v).FieldByName("MemberBase").Interface().(data.MemberBase)
	latestMsgID := r.getLatestMessage(roomBase.ID)

	field := make(map[string]interface{})
	field["member"] = memberBase.ID
	field["msgid"] = latestMsgID
	err := r.client.HMSet("room_"+roomBase.ID, field).Err()
	r.checkErr(err)
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

func (r *Redis) InsertRoom(room ...interface{}) {
	r.memoryStore.InsertRoom(room...)
}

func (r *Redis) GetRoom(rid ...string) []interface{} {
	return r.memoryStore.GetRoom(rid...)
}

func (r *Redis) InsertMember(member ...interface{}) {
	r.memoryStore.InsertMember(member...)
}

func (r *Redis) GetMember(mid ...string) []interface{} {
	return r.memoryStore.GetMember(mid...)
}

func (r *Redis) InsertToken(token string, mid string) {
	field := make(map[string]interface{})
	field["member"] = mid
	err := r.client.HMSet("token_"+token, field).Err()
	r.checkErr(err)
}

func (r *Redis) GetToken(token string) string {
	res, err := r.client.HMGet("token_"+token, "member").Result()
	r.checkErr(err)
	return res[0].(string)
}

func (r *Redis) InsertMessage(message *data.MessageBase) int {
	field := make(map[string]interface{})
	field["rid"] = message.RID
	field["mid"] = message.MID
	field["text"] = message.Text
	r.msgIDMutex.Lock()
	id, err := r.client.Incr("msg_id").Result()
	err = r.client.HMSet("room_"+message.RID+"_msg_"+strconv.Itoa(int(id)), field).Err()
	r.checkErr(err)

	err = r.client.LPush("room_"+message.RID+"_msg", strconv.Itoa(int(id))).Err()
	r.checkErr(err)

	r.msgIDMutex.Unlock()
	return int(id)
}

func (r *Redis) GetUnReadMessage(rid string, msg_id int) []*data.MessageBase {
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
	for i, msg := range res {
		if msg == strconv.Itoa(msg_id) {
			position = i
			break
		}
	}
	if len(res) == 0 {
		return messages
	}
	for _, msg := range res[:position] {
		message := &data.MessageBase{}
		res, err := r.client.HMGet("room_"+rid+"_msg_"+msg, "rid", "mid", "text").Result()
		r.checkErr(err)
		message.ID, _ = strconv.Atoi(msg)
		message.RID = res[0].(string)
		message.MID = res[1].(string)
		message.Text = res[2].(string)
		messages = append([]*data.MessageBase{message}, messages...)
	}
	return messages
}

func (r *Redis) GetHistoryMessage(rid string, msg_id int, offset int) []*data.MessageBase {
	sort := &redis.Sort{
		Order: "DESC",
	}
	position := -1
	start, end := 0, 0
	res, err := r.client.Sort("room_"+rid+"_msg", sort).Result()
	r.checkErr(err)
	for i, msg := range res {
		if msg == string(msg_id) {
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
	for _, msg := range res[start:end] {
		message := &data.MessageBase{}
		res, err := r.client.HMGet("room_"+rid+"_msg_"+msg, "rid", "mid", "text").Result()
		r.checkErr(err)
		message.RID = res[0].(string)
		message.MID = res[1].(string)
		message.Text = res[2].(string)
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

func (r *Redis) SetRoomMemberMessage(rid string, mid string, msg_id int) {
	field := make(map[string]interface{})
	field["msg_id"] = msg_id
	err := r.client.HMSet("room_"+rid+"_member_"+mid, field).Err()
	r.checkErr(err)
}

func (r *Redis) GetRoomMemberMessage(rid string, mid string) int {
	res, err := r.client.HMGet("room_"+rid+"_member_"+mid, "msg_id").Result()
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

func NewRedis(l log.Logger, config *config.RedisConfig) *Redis {
	redis := &Redis{}
	redis.initRedis(l, config)
	return redis
}
