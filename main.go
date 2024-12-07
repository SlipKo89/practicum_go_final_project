package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

const time_format = "20060102"

func fileServer(router *chi.Mux) {
	root := "web"
	fs := http.FileServer(http.Dir(root))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}

func main() {
	//webDir := "/Users/slipko/Desktop/Golang_learning/practicum_learning/final_task/go_final_project/web"
	//set default port number
	PortNum := "7540"
	//try get port number from env
	if os.Getenv("TODO_PORT") != "" {
		PortNum = os.Getenv("TODO_PORT")
	}

	exist_db()

	fmt.Println("Запускаем сервер на порту:", PortNum)
	PortNum = ":" + PortNum

	// test NextDate
	test, err := NextDate(time.Now(), "20241016", "y")
	fmt.Println(test, err)

	r := chi.NewRouter()
	fileServer(r)
	r.Post("/api/task", addTask)
	r.Get("/api/task", getTask)
	r.Put("/api/task", editTask)
	r.Get("/api/tasks", getTasks)
	r.Get("/api/nextdate", getNextDate)
	r.Post("/api/task/done", doneTask)
	r.Delete("/api/task", removeTask)
	http.ListenAndServe(PortNum, r)
	fmt.Println("Завершаем работу")

}
