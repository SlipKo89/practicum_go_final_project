package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

func editTask(w http.ResponseWriter, r *http.Request) {
	// create temporary variable for store data from UI interface after deserialize
	var task Task
	// create temporary variable for store raw data from UI interface before serialize
	var buf bytes.Buffer

	// read body from request, if it's empty - send error
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		createJsonError(w, "Ошибка чтения входных параметров")
		return
	}

	// deserialize raw data to json and store to var task
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		createJsonError(w, "Ошибка десериализации json")
		return
	}

	// check exist ID
	if task.ID == "" {
		createJsonError(w, "Не указан ID задачи")
		return
	} else {
		_, err = strconv.Atoi(task.ID)
		if err != nil {
			createJsonError(w, "ID задачи указан неверно")
			return
		}
	}

	// check exist title
	if task.Title == "" {
		createJsonError(w, "Не указан заголовок задачи")
		return
	}

	// check date format
	timeDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		createJsonError(w, "Неверный формат date")
		return
	}

	// check exist date
	if task.Date == "" {
		task.Date = time.Now().Format("20060201")
	}

	// set task date
	duration := timeDate.Sub(time.Now())
	if duration < 0 && task.Repeat == "" {
		task.Date = time.Now().Format("20060201")
	} else if duration < 0 && task.Repeat != "" {
		task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			createJsonError(w, "Неверный формат date")
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
	//fmt.Println("taskID", task.ID)
	_, err = db.Exec("UPDATE scheduler SET date = :date, comment = :comment, title = :title, repeat = :repeat WHERE id = :id",
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		createJsonError(w, "Невозможно добавить задачу в БД")
		//fmt.Println("Невозможно добавить задачу в БД", err)
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
