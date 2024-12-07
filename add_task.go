package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

func addTask(w http.ResponseWriter, r *http.Request) {
	// create temporary variable for store data from UI interface after deserialize
	var task Task
	// create temporary variable for store raw data from UI interface before serialize
	var buf bytes.Buffer

	t := time.Now()

	// read body from request, if it's empty - send error
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// deserialize raw data to json and store to var task
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		createJsonError(w, "Ошибка десериализации json")
		return
	}

	// check exist title
	if task.Title == "" {
		createJsonError(w, "Не указан заголовок задачи")
		return
	}

	// check exist date
	if task.Date == "" {
		task.Date = t.Format("20060102")
	}

	// check date format
	_, err = time.Parse("20060102", task.Date)
	if err != nil {
		createJsonError(w, "Неверный формат date")
		return
	}

	// check date duration
	taskDateInt, err := strconv.Atoi(task.Date)
	if err != nil {
		createJsonError(w, "Неверный формат date")
		return
	}

	nowDateInt, err := strconv.Atoi(time.Now().Format("20060102"))
	if err != nil {
		createJsonError(w, "Неверный формат date")
		return
	}

	if taskDateInt < nowDateInt && task.Repeat == "" {
		task.Date = t.Format("20060102")
		fmt.Println("Repeate 0", task.Date)
	} else if taskDateInt < nowDateInt && task.Repeat != "" {
		fmt.Println("Before NextDate", task.Date)
		task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		fmt.Println("After NextDate", task.Date)
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
	// create struct variable for create json
	type TaskRes struct {
		ID string `json:"id"`
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

/* тут была попытка сделать единую функцию для добавления и обновления задачи, но что-то пошло не так

func addTask(w http.ResponseWriter, r *http.Request) {
	// create temporary variable for store data from UI interface after deserialize
	var task Task
	// create temporary variable for store raw data from UI interface before serialize
	var buf bytes.Buffer

	t := time.Now()

	// read body from request, if it's empty - send error
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// deserialize raw data to json and store to var task
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		//createJsonError(w, errors.New("Ошибка десериализации json"))
		createJsonError(w, "Ошибка десериализации json")
		return
	}

	// check exist title
	if task.Title == "" {
		//createJsonError(w, errors.New("Не указан заголовок задачи"))
		createJsonError(w, "Не указан заголовок задачи")
		return
	}

	// check exist date
	if task.Date == "" {
		task.Date = t.Format("20060102")
	}

	// check date format
	_, err = time.Parse("20060102", task.Date)
	if err != nil {
		//createJsonError(w, errors.New("Неверный формат date"))
		createJsonError(w, "Неверный формат date")
		return
	}

	// check date duration
	taskDateInt, err := strconv.Atoi(task.Date)
	if err != nil {
		createJsonError(w, "Неверный формат date")
		return
	}

	nowDateInt, err := strconv.Atoi(time.Now().Format("20060102"))
	if err != nil {
		createJsonError(w, "Неверный формат date")
		return
	}

	if taskDateInt < nowDateInt && task.Repeat == "" {
		task.Date = t.Format("20060102")
		fmt.Println("Repeate 0", task.Date)
	} else if taskDateInt < nowDateInt && task.Repeat != "" {
		fmt.Println("Before NextDate", task.Date)
		task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		fmt.Println("After NextDate", task.Date)
		if err != nil {
			//createJsonError(w, errors.New("Неверный формат date"))
			createJsonError(w, "Неверный формат date")
			return
		}
	}

	// open db
	db, err := sql.Open("sqlite", DBFilePath)
	if err != nil {
		//createJsonError(w, errors.New("Невозможно открыть файл БД"))
		createJsonError(w, "Невозможно открыть файл БД")
		return
	}
	// auto close after finish all operations
	defer db.Close()

	if task.ID == "" {
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

		// create struct variable for create json
		type TaskRes struct {
			ID string `json:"id",omitempty`
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

	} else {
		// modify task in db
		_, err := db.Exec("UPDATE scheduler SET date = :date, comment = :comment, title = :title, repeat = :repeat WHERE id = :id",
			sql.Named("id", task.ID),
			sql.Named("date", task.Date),
			sql.Named("title", task.Title),
			sql.Named("comment", task.Comment),
			sql.Named("repeat", task.Repeat))
		if err != nil {
			createJsonError(w, "Невозможно обновить задачу в БД")
			fmt.Println(err)
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
}
*/
