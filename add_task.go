package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite"
)

func addTask(w http.ResponseWriter, r *http.Request) {

	// create temporary variable for store raw data from UI interface before serialize
	var buf bytes.Buffer

	// read body from request, if it's empty - send error
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		createJsonError(w, "Ошибка чтения входных параметров")
		return
	}

	// use shared func for addTask and editTask
	error_string := add_edit_task(buf)
	if error_string != "" {
		createJsonError(w, error_string)
		return
	}

	// open db
	db, err := sql.Open("sqlite", DBFilePath)
	if err != nil {
		createJsonError(w, "Невозможно открыть файл БД")
		return
	}
	// auto close after finish all operations
	defer db.Close()

	// add new task in db
	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		createJsonError(w, "Невозможно добавить задачу в БД")
		return
	}

	// get last added row id
	id, err := res.LastInsertId()
	if err != nil {
		createJsonError(w, "Невозможно получить id последней добавленной записи")
		return
	}

	var TaskResult = TaskRes{
		ID: strconv.FormatInt(id, 10),
	}

	// сериализация ответа
	resp, err := json.Marshal(TaskResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		createJsonError(w, "Не прошла сериализация данных")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)

}
