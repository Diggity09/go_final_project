package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
	"strconv"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var apiTask Task

	// Десериализуем JSON
	if err := json.NewDecoder(r.Body).Decode(&apiTask); err != nil {
		writeJSON(w, http.StatusBadRequest, UpdateTaskResponse{Error: "Ошибка десериализации JSON"})
		return
	}

	// Проверяем, что указан ID
	if apiTask.ID == "" {
		writeJSON(w, http.StatusBadRequest, UpdateTaskResponse{Error: "Не указан идентификатор задачи"})
		return
	}

	// Проверяем обязательное поле title
	if apiTask.Title == "" {
		writeJSON(w, http.StatusBadRequest, UpdateTaskResponse{Error: "Не указан заголовок задачи"})
		return
	}

	// Проверяем и корректируем дату
	if err := checkDate(&apiTask); err != nil {
		writeJSON(w, http.StatusBadRequest, UpdateTaskResponse{Error: err.Error()})
		return
	}

	// Преобразуем строковый ID в int64
	idInt, err := strconv.ParseInt(apiTask.ID, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, UpdateTaskResponse{Error: "Некорректный идентификатор"})
		return
	}

	// Преобразуем в структуру базы данных
	dbTask := &db.Task{
		ID:      idInt,
		Date:    apiTask.Date,
		Title:   apiTask.Title,
		Comment: apiTask.Comment,
		Repeat:  apiTask.Repeat,
	}

	// Обновляем задачу в базе данных
	if err := db.UpdateTask(dbTask); err != nil {
		if err.Error() == "incorrect id for updating task" {
			writeJSON(w, http.StatusBadRequest, UpdateTaskResponse{Error: "Задача не найдена"})
		} else {
			writeJSON(w, http.StatusInternalServerError, UpdateTaskResponse{Error: "Ошибка обновления задачи в базе данных"})
		}
		return
	}

	// Возвращаем пустой JSON в случае успеха
	writeJSON(w, http.StatusOK, map[string]interface{}{})
}
