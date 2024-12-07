package main

import (
	"encoding/json"
	"net/http"

	_ "modernc.org/sqlite"
)

func removeTask(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	// check exist task
	error_status := check_exist_task(id)
	if error_status != "" {
		createJsonError(w, error_status)
		return
	}

	error_status = delete_task_from_db(id)
	if error_status != "" {
		createJsonError(w, error_status)
		return
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
