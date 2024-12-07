package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	_ "modernc.org/sqlite"
)

func getTasks(w http.ResponseWriter, r *http.Request) {

	tasks := make([]Task, 0)

	db, err := sql.Open("sqlite", DBFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			fmt.Println(err)
			return
		}

		tasks = append(tasks, task)
	}

	// Sort tasks in slice
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Date < tasks[j].Date
	})

	tasksMap := make(map[string][]Task)
	tasksMap["tasks"] = append(tasks)

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return
	}

	// сериализуем данные из слайса artists
	resp, err := json.Marshal(tasksMap)
	if err != nil {
		createJsonError(w, "Ошибка сериализации")
		return
	}

	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

func getTask(w http.ResponseWriter, r *http.Request) {
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
	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		createJsonError(w, "Задача не найдена")
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		createJsonError(w, "Не прошла сериализация данных")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
