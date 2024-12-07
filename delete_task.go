package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "modernc.org/sqlite"
)

func removeTask(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	var task Task
	// open connect to DB
	db, err := sql.Open("sqlite", DBFilePath)
	if err != nil {
		createJsonError(w, "Файл БД открыть не получается")
		return
	}
	defer db.Close()

	// check exist id
	row := db.QueryRow("SELECT id FROM scheduler WHERE id = :id", sql.Named("id", id))
	err = row.Scan(&task.ID)
	if err != nil {
		createJsonError(w, "Задача не найдена")
		return
	}

	// delete task from DB
	_, err = db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		createJsonError(w, "Ошибка удаления записи")
		return
	}

	var emptyStruct = empty{}
	resp, err := json.Marshal(emptyStruct)
	if err != nil {
		createJsonError(w, "Ошибка сериализации")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
