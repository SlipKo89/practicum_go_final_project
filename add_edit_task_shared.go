package main

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

// create temporary variable for store data from UI interface after deserialize
var task Task

func add_edit_task(buf bytes.Buffer) string {

	t := time.Now()

	// deserialize raw data to json and store to var task
	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		return "Ошибка десериализации json"
	}

	// check exist title
	if task.Title == "" {
		return "Не указан заголовок задачи"
	}

	// check exist date
	if task.Date == "" {
		task.Date = time.Now().Format(time_format)
	} else {

		// check date format
		_, err := time.Parse(time_format, task.Date)
		if err != nil {
			return "Неверный формат date. Parse error"
		}
	}

	// check date duration
	taskDateInt, err := strconv.Atoi(task.Date)
	if err != nil {
		return "Неверный формат date"
	}

	nowDateInt, err := strconv.Atoi(time.Now().Format(time_format))
	if err != nil {
		return "Неверный формат текущего времени"
	}

	if taskDateInt < nowDateInt && task.Repeat == "" {
		task.Date = t.Format(time_format)
	} else if taskDateInt < nowDateInt && task.Repeat != "" {
		task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return "Ошибка Nextdate"
		}
	}

	return ""
}
