package main

import (
	"encoding/json"
	"net/http"
	"fmt"

	_ "modernc.org/sqlite"
)

//type AddTaskError struct {
//	Error error `json:"error"`
//}

type AddTaskErrorS struct {
	Error string `json:"error"`
}


func createJsonError(w http.ResponseWriter, error_name string) {
	var ErrorResult = AddTaskErrorS{
		Error: error_name,
	}

	fmt.Println(ErrorResult.Error)

	resp, err := json.Marshal(ErrorResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(resp)
}
