package api

import (
	"go_final_project/pkg/db"
	"net/http"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр id из URL
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Не указан идентификатор"})
		return
	}

	// Удаляем задачу из базы данных
	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Задача не найдена"})
		return
	}

	// Возвращаем пустой JSON в случае успеха
	writeJSON(w, http.StatusOK, map[string]interface{}{})
}
