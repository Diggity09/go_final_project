package api

import (
	"go_final_project/pkg/db"
	"net/http"
	"strconv"
)

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	dbTasks, err := db.Tasks(50)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Ошибка получения задач из базы данных"})
		return
	}

	// Преобразуем в API структуры
	var apiTasks []*Task
	for _, dbTask := range dbTasks {
		if dbTask != nil { // Проверяем на nil
			apiTask := &Task{
				ID:      strconv.FormatInt(dbTask.ID, 10),
				Date:    dbTask.Date,
				Title:   dbTask.Title,
				Comment: dbTask.Comment,
				Repeat:  dbTask.Repeat,
			}
			apiTasks = append(apiTasks, apiTask)
		}
	}

	// Если нет задач, возвращаем пустой массив
	if apiTasks == nil {
		apiTasks = []*Task{}
	}

	writeJSON(w, http.StatusOK, TasksResponse{Tasks: apiTasks})
}
