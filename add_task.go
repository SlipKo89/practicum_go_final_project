package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type Task struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func addTask(w http.ResponseWriter, r *http.Request) {
	// create temporary variable for store data from UI interface after deserialize
	var task Task
	// create temporary variable for store raw data from UI interface before serialize
	var buf bytes.Buffer

	// read body from request, if it's empty - send error
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// deserialize raw data to json and store to var task
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// check empty title
	if task.Title == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// check date format
	_, err = time.Parse("20060102", task.Date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// temporary print var task
	fmt.Println(task)

	// open db
	db, err := sql.Open("sqlite", DBFilePath)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		return
	}

	id, err := res.LastInsertId()

	type AddTaskResult struct {
		ID int64 `json:"id"`
	}

	var TaskResult = AddTaskResult{
		ID: id,
	}

	resp, err := json.Marshal(TaskResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err != nil {
		type AddTaskError struct {
			Error error `json:"error"`
		}
		var ErrorResult = AddTaskError{
			Error: err,
		}
		resp, err = json.Marshal(ErrorResult)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
