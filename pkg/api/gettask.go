package api

import (
	"go_final_project/pkg/db"
	"net/http"
	"strconv"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр id из URL
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Не указан идентификатор"})
		return
	}

	// Получаем задачу из базы данных
	dbTask, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "Задача не найдена"})
		return
	}

	// Преобразуем в API структуру
	apiTask := Task{
		ID:      strconv.FormatInt(dbTask.ID, 10),
		Date:    dbTask.Date,
		Title:   dbTask.Title,
		Comment: dbTask.Comment,
		Repeat:  dbTask.Repeat,
	}

	// Возвращаем задачу
	writeJSON(w, http.StatusOK, apiTask)
}
