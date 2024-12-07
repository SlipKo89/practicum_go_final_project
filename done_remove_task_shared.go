package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func check_exist_task(id string) string {

	// open connect to DB
	db, err := sql.Open("sqlite", DBFilePath)
	if err != nil {
		return "Файл БД открыть не получается"
	}
	defer db.Close()

	// check exist id
	row := db.QueryRow("SELECT id, date, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))
	err = row.Scan(&task.ID, &task.Date, &task.Repeat)
	if err != nil {
		return "Задача не найдена"
	}

	return ""
}

func delete_task_from_db(id string) string {

	// open connect to DB
	db, err := sql.Open("sqlite", DBFilePath)
	if err != nil {
		return "Файл БД открыть не получается"
	}
	defer db.Close()

	// delete task from DB
	_, err = db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return "Ошибка удаления задачи из БД"
	}

	return ""
}
