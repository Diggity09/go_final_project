package server

import (
	"net/http"
)

func Run(port string) error {
	// Обслуживание статических файлов
	http.Handle("/", http.FileServer(http.Dir("./web")))

	return http.ListenAndServe(":"+port, nil)
}
