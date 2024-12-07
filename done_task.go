package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

func doneTask(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	var task Task

	// open db
	db, err := sql.Open("sqlite", DBFilePath)
	if err != nil {
		createJsonError(w, "Файл БД открыть не получается")
		return
	}
	defer db.Close()

	// make request to db for find task with id
	row := db.QueryRow("SELECT id, date, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))
	err = row.Scan(&task.ID, &task.Date, &task.Repeat)
	if err != nil {
		createJsonError(w, "Задача не найдена")
		return
	}

	// check repeat
	if task.Repeat == "" {
		_, err = db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
		if err != nil {
			createJsonError(w, "Ошибка удаления записи")
			return
		}
	} else {
		// получаем новую дату задачи при условии непустого repeat
		resp, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			createJsonError(w, "Невозможно вычислить следующую дату")
			return
		}
		// add new task in db
		_, err = db.Exec("UPDATE scheduler SET date = :date WHERE id = :id",
			sql.Named("id", task.ID),
			sql.Named("date", resp))
		if err != nil {
			createJsonError(w, "Невозможно обновить задачу в БД")
			return
		}
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
