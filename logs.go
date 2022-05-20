package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

func NewCommandLog(from, command, result string) *CommandLog {
	return &CommandLog{From: from, Command: command, Result: result, Time: int(time.Now().Unix())}
}

type CommandLog struct {
	Id int
	From string
	Command string
	Result string
	Time int
}

func (c *CommandLog) log() error {
	_, err := db.Exec("INSERT INTO `commands` (`from`, `command`, `result`, `time`)" +
		"VALUES (?, ?, ?, ?)",
		c.From, c.Command, c.Result, c.Time)
	return err
}

var db *sql.DB

func initLogs() (err error) {
	db, err = sql.Open("mysql", "root:root@/msu_go_11?charset=utf8&interpolateParams=true")
	if err != nil {
		return
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `commands` (" +
		"`id` serial NOT NULL," +
		"`from` text NOT NULL," +
		"`command` text NOT NULL, " +
		"`result` text NOT NULL, " +
		"`time` integer NOT NULL" +
		");")
	return
}

func getLogs(n uint) (logs []CommandLog) {
	var log CommandLog
	rows, err := db.Query("SELECT * FROM `commands` LIMIT " + strconv.FormatUint(uint64(n), 10) + ";")
	panicOnError(err)
	for rows.Next() {
		var id int
		var from string
		var command string
		var result string
		var timestamp int
		panicOnError(rows.Scan(&id, &from, &command, &result, &timestamp))
		log = CommandLog{id, from, command, result, timestamp}
		logs = append(logs, log)
	}
	return
}
