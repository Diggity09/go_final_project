package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
	"strconv"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var apiTask Task

	// Десериализуем JSON
	if err := json.NewDecoder(r.Body).Decode(&apiTask); err != nil {
		writeJSON(w, http.StatusBadRequest, AddTaskResponse{Error: "Ошибка десериализации JSON"})
		return
	}

	// Проверяем обязательное поле title
	if apiTask.Title == "" {
		writeJSON(w, http.StatusBadRequest, AddTaskResponse{Error: "Не указан заголовок задачи"})
		return
	}

	// Проверяем и корректируем дату
	if err := checkDate(&apiTask); err != nil {
		writeJSON(w, http.StatusBadRequest, AddTaskResponse{Error: err.Error()})
		return
	}

	// Преобразуем в структуру базы данных
	dbTask := &db.Task{
		Date:    apiTask.Date,
		Title:   apiTask.Title,
		Comment: apiTask.Comment,
		Repeat:  apiTask.Repeat,
	}

	// Добавляем задачу в базу данных
	id, err := db.AddTask(dbTask)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, AddTaskResponse{Error: "Ошибка добавления задачи в базу данных"})
		return
	}

	// Возвращаем успешный ответ с правильным ID
	writeJSON(w, http.StatusOK, AddTaskResponse{ID: strconv.FormatInt(id, 10)})
}
