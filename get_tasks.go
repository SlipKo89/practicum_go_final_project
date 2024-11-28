package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "modernc.org/sqlite"
)

type TaskWithID struct {
	Task
	ID string `json:"id"`
}

func getTasks(w http.ResponseWriter, r *http.Request) {

	tasks := make([]TaskWithID, 0)

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
		task := TaskWithID{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(task)
		tasks = append(tasks, task)
	}

	tasksMap := make(map[string][]TaskWithID)
	tasksMap["tasks"] = append(tasks)

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return
	}

	// сериализуем данные из слайса artists
	resp, err := json.Marshal(tasksMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(resp)

	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}
