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

	// check exist id
	error_status := check_exist_task(id)
	if error_status != "" {
		createJsonError(w, error_status)
		return
	}

	// check repeat
	if task.Repeat == "" {
		error_status = delete_task_from_db(id)
		if error_status != "" {
			createJsonError(w, error_status)
			return
		}

	} else {
		// получаем новую дату задачи при условии непустого repeat
		resp, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			createJsonError(w, "Невозможно вычислить следующую дату")
			return
		}

		// open db
		db, err := sql.Open("sqlite", DBFilePath)
		if err != nil {
			createJsonError(w, "Файл БД открыть не получается")
			return
		}
		defer db.Close()

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
