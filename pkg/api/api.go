package api

import (
	"encoding/json"
	"log"
	"net/http"
)

const DateFormat = "20060102"

func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", doneHandler)
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	writeJson(w, map[string]string{"error": message})
}

func writeErrorWithLog(w http.ResponseWriter, message string, err error, statusCode int) {
	if err != nil {
		log.Printf("Ошибка: %s - %v", message, err)
	}
	writeError(w, message, statusCode)
}
