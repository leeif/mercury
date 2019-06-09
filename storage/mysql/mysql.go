package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"

	avl "github.com/Workiva/go-datastructures/tree/avl"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/leeif/mercury/storage/config"
	"github.com/leeif/mercury/storage/data"
	"github.com/leeif/mercury/storage/memory"
)

const createMercury = "create database if not exists `mercury`"

const createRoom = "create table if not exists mercury.`room` (" +
	"`id`    int(10) not null auto_increment," +
	// room id
	"`rid`   varchar(100) not null," +
	// mid id
	"`mid`   varchar(100) not null," +
	// latest read msgid
	"`msgid` int(10) not null," +
	"primary key(`id`)," +
	"unique  key(`rid`, `mid`)," +
	"unique  key(`rid`, `mid`, `msgid`)," +
	"key(`mid`, `rid`)" +
	")"

const createMessage = "create table if not exists mercury.`message` (" +
	// message id
	"`id`    int(10) not null auto_increment," +
	// room id
	"`rid`   varchar(100) not null," +
	// member id
	"`mid`   varchar(100) not null," +
	"`text`  text not null," +
	// create time
	"`ctime` timestamp not null," +
	"primary key(`id`)" +
	")"

const createToken = "create table if not exists mercury.`token` (" +
	"`id`    int(10) not null auto_increment," +
	// member id
	"`mid`   varchar(100) not null," +
	// token
	"`token` varchar(100) not null," +
	// create time
	"`ctime` timestamp not null DEFAULT CURRENT_TIMESTAMP," +
	"primary key(`id`)" +
	")"

type MySQL struct {
	logger      log.Logger
	db          *sql.DB
	memoryStore *memory.Memory
}

func (m *MySQL) initDB(l log.Logger, config *config.MySQLConfig) {
	var err error
	m.logger = l
	conString := config.User + ":" + config.Password +
		"@tcp(" + config.Host + ":" + config.Port + ")/?parseTime=true"
	level.Debug(m.logger).Log("connection", conString)
	m.db, err = sql.Open("mysql", conString)
	m.fatalErr(err)

	m.db.SetMaxOpenConns(20)
	m.db.SetMaxIdleConns(10)

	m.memoryStore = &memory.Memory{
		Room:   avl.NewImmutable(),
		Member: avl.NewImmutable(),
	}

	initQuery := []string{
		createMercury,
		createRoom,
		createMessage,
		createToken,
	}
	for _, query := range initQuery {
		var err error
		var res sql.Result
		res, err = m.db.Exec(query)
		m.fatalErr(err)
		_, err = res.RowsAffected()
		m.fatalErr(err)
	}
}

func (m *MySQL) InsertRoomMember(room interface{}, member interface{}) {
	var v reflect.Value
	var query string
	v = reflect.ValueOf(room)
	roomBase := reflect.Indirect(v).FieldByName("RoomBase").Interface().(data.RoomBase)
	v = reflect.ValueOf(member)
	memberBase := reflect.Indirect(v).FieldByName("MemberBase").Interface().(data.MemberBase)
	latestMsgID := m.getLatestMessage(roomBase.ID)
	if m.checkIfExistRoomMember(roomBase.ID, memberBase.ID) {
		return
	}

	query = "insert into mercury.`room` (`rid`, `mid`, `msgid`) value(?, ?, ?)"
	stmt, err := m.db.Prepare(query)
	m.checkErr(err)
	res, err := stmt.Exec(roomBase.ID, memberBase.ID, latestMsgID)
	m.checkErr(err)
	id, err := res.LastInsertId()
	fmt.Printf(string(id))
	m.checkErr(err)
}

func (m *MySQL) checkIfExistRoomMember(rid string, mid string) bool {
	query := "select count(*) from mercury.`room` where rid=? and mid=?"
	rows, err := m.db.Query(query, rid, mid)
	m.checkErr(err)
	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		m.checkErr(err)
	}
	level.Debug(m.logger).Log("count", count)
	if count > 0 {
		return true
	}
	return false
}

func (m *MySQL) getLatestMessage(rid string) int {
	query := "select `id` from mercury.`message` where rid=? order by `id` desc limit 1"
	rows, err := m.db.Query(query, rid)
	m.checkErr(err)
	var msgid int
	for rows.Next() {
		err := rows.Scan(&msgid)
		m.checkErr(err)
	}
	return msgid
}

func (m *MySQL) InsertRoom(room ...interface{}) {
	m.memoryStore.InsertRoom(room...)
}

func (m *MySQL) GetRoom(rid ...string) []interface{} {
	return m.memoryStore.GetRoom(rid...)
}

func (m *MySQL) InsertMember(member ...interface{}) {
	m.memoryStore.InsertMember(member...)
}

func (m *MySQL) GetMember(mid ...string) []interface{} {
	return m.memoryStore.GetMember(mid...)
}

func (m *MySQL) InsertToken(token string, mid string) {
	var query string
	var rows *sql.Rows
	var res sql.Result
	var err error
	var stmt *sql.Stmt
	query = "select token from mercury.`token` where mid=?"
	rows, err = m.db.Query(query, mid)
	m.checkErr(err)
	exists := false
	for rows.Next() {
		exists = true
	}
	if !exists {
		query = "insert into mercury.`token` (mid, token) value (?, ?)"
		stmt, err = m.db.Prepare(query)
		m.checkErr(err)
		res, err = stmt.Exec(mid, token)
	} else {
		query = "update mercury.`token` set token=? where mid=?"
		stmt, err = m.db.Prepare(query)
		m.checkErr(err)
		res, err = stmt.Exec(token, mid)
	}
	m.checkErr(err)
	id, err := res.LastInsertId()
	fmt.Printf(string(id))
	m.checkErr(err)
}

func (m *MySQL) GetToken(token string) string {
	query := "select mid from mercury.`token` where token=?"
	rows, err := m.db.Query(query, token)
	m.checkErr(err)
	var mid string
	for rows.Next() {
		err := rows.Scan(&mid)
		m.checkErr(err)
	}
	return mid
}

func (m *MySQL) InsertMessage(message *data.MessageBase) int {
	query := "insert into mercury.`message` (`rid`, `mid`, `text`) value(?, ?, ?)"
	stmt, err := m.db.Prepare(query)
	m.checkErr(err)
	res, err := stmt.Exec(message.RID, message.MID, message.Text)
	m.checkErr(err)
	id, err := res.LastInsertId()
	fmt.Printf(string(id))
	m.checkErr(err)
	level.Debug(m.logger).Log(id)
	return int(id)
}

func (m *MySQL) GetUnReadMessage(rid string, msg_id int) []*data.MessageBase {
	query := "select `id`, `rid`, `mid`, `text`, `ctime` from mercury.`message` where rid=? and id > ?"
	rows, err := m.db.Query(query, rid, msg_id)
	m.checkErr(err)
	messages := make([]*data.MessageBase, 0)
	for rows.Next() {
		message := &data.MessageBase{}
		err := rows.Scan(&message.ID, &message.RID, &message.MID, &message.Text, &message.CTime)
		m.checkErr(err)
		messages = append(messages, message)
	}
	return messages
}

func (m *MySQL) GetHistoryMessage(rid string, msg_id int, offset int) []*data.MessageBase {
	query := "select `id`, `rid`, `mid`, `text`, `ctime` from mercury.`message` where rid=? and id < ? limit ?"
	rows, err := m.db.Query(query, rid, msg_id, offset)
	m.checkErr(err)
	messages := make([]*data.MessageBase, 0)
	for rows.Next() {
		message := data.MessageBase{}
		err := rows.Scan(&message.ID, &message.RID, &message.MID, &message.Text, &message.CTime)
		m.checkErr(err)
		messages = append(messages, &message)
	}
	return messages
}

func (m *MySQL) SetMemberOfRoom(rid string, mid string) {
	// do nothing
	return
}

func (m *MySQL) GetMemberFromRoom(rid string) []string {
	query := "select `mid` from mercury.room where rid=?"
	rows, err := m.db.Query(query, rid)
	m.checkErr(err)
	res := make([]string, 0)
	for rows.Next() {
		var mid string
		err := rows.Scan(&mid)
		m.checkErr(err)
		res = append(res, mid)
	}
	return res
}

func (m *MySQL) SetRoomOfMember(mid string, rid string) {
	// do nothing
	return
}

func (m *MySQL) GetRoomFromMember(mid string) []string {
	query := "select `rid` from mercury.message where `mid`=?"
	rows, err := m.db.Query(query, mid)
	m.checkErr(err)
	res := make([]string, 0)
	for rows.Next() {
		var rid string
		err := rows.Scan(&rid)
		m.checkErr(err)
		res = append(res, rid)
	}
	return res
}

func (m *MySQL) SetRoomMemberMessage(rid string, mid string, msg_id int) {
	query := "update mercury.`room` set msgid=? where rid=? and mid=?"
	stmt, err := m.db.Prepare(query)
	m.checkErr(err)
	res, err := stmt.Exec(msg_id, rid, mid)
	m.checkErr(err)
	id, err := res.RowsAffected()
	fmt.Printf(string(id))
	m.checkErr(err)
}

func (m *MySQL) GetRoomMemberMessage(rid string, mid string) int {
	query := "select `msgid` from mercury.`room` where rid=? and mid=?"
	rows, err := m.db.Query(query, rid, mid)
	m.checkErr(err)
	var msgid int
	for rows.Next() {
		err := rows.Scan(&msgid)
		m.checkErr(err)
	}
	return msgid
}

func (m *MySQL) checkErr(err error) {
	if err != nil {
		level.Error(m.logger).Log("msg", err)
	}
}

func (m *MySQL) fatalErr(err error) {
	if err != nil {
		level.Error(m.logger).Log("msg", err)
		os.Exit(1)
	}
}

func NewMySQL(l log.Logger, config *config.MySQLConfig) *MySQL {
	mysql := &MySQL{}
	mysql.initDB(l, config)
	return mysql
}
