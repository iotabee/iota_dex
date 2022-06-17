package model

import (
	"database/sql"
	"fmt"
	"iota_dex/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func ConnectToMysql() {
	var err error
	db, err = sql.Open("mysql", config.Db.Usr+":"+config.Db.Pwd+"@tcp("+config.Db.Host+":"+config.Db.Port+")/"+config.Db.DbName)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); nil != err {
		log.Panic("Connect to Mysql error : " + err.Error())
	}
}

func Ping() error {
	if db == nil {
		return fmt.Errorf("mysql connection is nil")
	}
	if err := db.Ping(); nil != err {
		return fmt.Errorf("connect to Mysql error : %v", err)
	}
	return nil
}
