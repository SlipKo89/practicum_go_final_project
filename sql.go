package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DBFilePath = "scheduler.db"

func exist_db() {
	//try get DBFile path from env
	if os.Getenv("TODO_DBFILE") != "" {
		DBFilePath = os.Getenv("TODO_DBFILE")
	}

	//check db file exist
	_, err := os.Stat(DBFilePath)

	var install bool
	if err != nil {
		install = true
	}

	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	if install {
		db, err := sql.Open("sqlite", DBFilePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer db.Close()

		_, err1 := db.Exec("CREATE TABLE IF NOT EXISTS scheduler  (id	INTEGER NOT NULL UNIQUE, date INTEGER, title TEXT, comment	TEXT, repeat	TEXT, PRIMARY KEY(id AUTOINCREMENT))")

		if err1 != nil {
			fmt.Println(err1)
			return
		}

		_, err2 := db.Exec("CREATE INDEX date ON scheduler (date);")

		if err2 != nil {
			fmt.Println(err2)
			return
		}
	}
}
