package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite"
)

func editTask(w http.ResponseWriter, r *http.Request) {

	// create temporary variable for store raw data from UI interface before serialize
	var buf bytes.Buffer

	// read body from request, if it's empty - send error
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		createJsonError(w, "Ошибка чтения входных параметров")
		return
	}

	error_string := add_edit_task(buf)
	if error_string != "" {
		createJsonError(w, "Общие проверки не пройдены")
		return
	}

	// check exist ID
	if task.ID == "" {
		createJsonError(w, "Не указан ID задачи")
		return
	} else {
		_, err := strconv.Atoi(task.ID)
		if err != nil {
			createJsonError(w, "ID задачи указан неверно")
			return
		}
	}

	// open db
	db, err := sql.Open("sqlite", DBFilePath)
	if err != nil {
		createJsonError(w, "Невозможно открыть файл БД")
		return
	}
	// auto close after finish all operations
	defer db.Close()

	// check exist id
	row := db.QueryRow("SELECT id FROM scheduler WHERE id = :id", sql.Named("id", task.ID))
	err = row.Scan(&task.ID)
	if err != nil {
		createJsonError(w, "Задача не найдена")
		return
	}
	// add new task in db
	_, err = db.Exec("UPDATE scheduler SET date = :date, comment = :comment, title = :title, repeat = :repeat WHERE id = :id",
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		createJsonError(w, "Невозможно добавить задачу в БД")
		return
	}

	var emptyStruct = empty{}
	resp, err := json.Marshal(emptyStruct)
	if err != nil {
		createJsonError(w, "Ошибка сериализации")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
