package main

import (
	"encoding/json"
	"net/http"

	_ "modernc.org/sqlite"
)

func createJsonError(w http.ResponseWriter, error_name string) {
	var ErrorResult = AddTaskError{
		Error: error_name,
	}

	resp, err := json.Marshal(ErrorResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(resp)
}
