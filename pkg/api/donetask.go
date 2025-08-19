package api

import (
	"fmt"
	"go_final_project/pkg/db"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Задача не найдена"})
		return
	}

	// Добавляем логирование для отладки
	fmt.Printf("DEBUG: Task ID: %s, Current Date: %s, Repeat: %s\n", id, task.Date, task.Repeat)

	if task.Repeat != "" {
		taskDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Ошибка парсинга даты задачи"})
			return
		}

		parts := strings.Fields(task.Repeat)
		if len(parts) == 0 {
			if err := db.DeleteTask(id); err != nil {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Ошибка удаления задачи"})
				return
			}
			writeJSON(w, http.StatusOK, map[string]interface{}{})
			return
		}

		var nextDateStr string
		var nextErr error

		if len(parts) == 2 && parts[0] == "d" {
			interval, err := strconv.Atoi(parts[1])
			if err != nil || interval <= 0 || interval > 400 {
				if err := db.DeleteTask(id); err != nil {
					writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Ошибка удаления задачи"})
					return
				}
				writeJSON(w, http.StatusOK, map[string]interface{}{})
				return
			}
			nextDate := taskDate.AddDate(0, 0, interval)
			nextDateStr = nextDate.Format("20060102")

			// Добавляем логирование
			fmt.Printf("DEBUG: Adding %d days to %s = %s\n", interval, task.Date, nextDateStr)
		} else {
			nextDateStr, nextErr = NextDate(time.Now(), task.Date, task.Repeat)
		}

		if nextErr != nil {
			if err := db.DeleteTask(id); err != nil {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Ошибка удаления задачи"})
				return
			}
		} else {
			if err := db.UpdateTaskDate(nextDateStr, id); err != nil {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Ошибка обновления даты задачи"})
				return
			}
		}
	} else {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Задача не найдена"})
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{})
}
