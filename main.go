package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

/*
	func NextDate(now time.Time, date string, repeat string) (string, error) {
		if repeat == "" {
			return "", errors.New("repeat field is empty")
		}

		repeatParams := strings.Split(repeat, " ")

		if len(repeatParams) == 1 && repeatParams[0] == "y" {

		}


		if len(repeatParams) == 2 && repeatParams[0] == "d" && 0 < strconv.Atoi(repeatParams[1]) <= 400 {

		}

}
*/

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

	r := chi.NewRouter()
	fileServer(r)
	r.Post("/api/task", addTask)
	r.Get("/api/task", getTasks)
	http.ListenAndServe(PortNum, r)
	fmt.Println("Завершаем работу")
}
