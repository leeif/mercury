package mysql

import (
	"database/sql"
	"github.com/leeif/mercury/storage/config"
	"github.com/leeif/mercury/storage/data"
	"fmt"
)

const (
	DATABASE     = "mercury"
	TABLETOKEN   = "token"
	TABLEMESSAGE = "message"
	TABLEROOM    = "room"
)

var (
	db *sql.DB
)

func initDB(config *config.StorageConfig) {
	var err error
	conString := config.MySQLConfig.User + ":" + config.MySQLConfig.Password + 
		"@tcp( " + config.MySQLConfig.Host + ":" + config.MySQLConfig.Password + " )"
	db, err = sql.Open("mysql", conString)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	checkErr(err)
	createDatabase()
	createTable()
}

func createDatabase() {

}

func createTable() {
}

type MemberInMySQL struct {
}

func (m *MemberInMySQL) Insert(members data.MemberBase) {
	query := ""
	stmt, err := db.Prepare(query)
	checkErr(err)
	res, err := stmt.Exec()
	checkErr(err)
	id, err := res.LastInsertId()
	fmt.Printf(string(id))
	checkErr(err)
}

func (m *MemberInMySQL) Get(mid ...string) []data.MemberBase {
	query := ""
	rows, err := db.Query(query)
	if err != nil {

	}
	defer rows.Close()
	res := make([]data.MemberBase, 0)
	return res
}

type RoomInMySQL struct {
}

type MessageInMySQL struct {
}

type TokenInMySQL struct {
}

type IndexInMySQL struct {
}


func checkErr(err error) {

}
