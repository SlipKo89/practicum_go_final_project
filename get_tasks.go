package main

import (
	"encoding/json"
	"net/http"

	_ "modernc.org/sqlite"
)

func getTasks(w http.ResponseWriter, r *http.Request) {
	// сериализуем данные из слайса artists
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}
