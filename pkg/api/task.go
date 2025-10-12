package api

import (
	"encoding/json"
	"fmt"
	"go_final_project/pkg/db"
	"net/http"
	"strconv"
	"time"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeErrorWithLog(w, "Ошибка добавления задачи", err, http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	writeJson(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeError(w, "Не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	writeJson(w, map[string]interface{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeError(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	writeJson(w, map[string]interface{}{})
}

func doneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		writeError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	// Если нет правила повторения, удаляем задачу
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeErrorWithLog(w, "Ошибка удаления задачи", err, http.StatusInternalServerError)
			return
		}
	} else {
		// Вычисляем следующую дату
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeError(w, "Ошибка вычисления следующей даты", http.StatusBadRequest)
			return
		}

		// Обновляем дату задачи
		if err := db.UpdateDate(nextDate, id); err != nil {
			writeErrorWithLog(w, "Ошибка обновления даты задачи", err, http.StatusInternalServerError)
			return
		}
	}

	writeJson(w, map[string]interface{}{})
}

func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Truncate(24 * time.Hour)

	// Если дата не указана, используем сегодняшнюю
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
		return nil
	}

	// Проверяем корректность даты
	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("дата представлена в некорректном формате")
	}

	// Если есть правило повторения, проверяем его корректность
	if task.Repeat != "" {
		_, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("правило повторения указано в неправильном формате")
		}
	}

	// Если дата меньше сегодняшней (но не равна ей)
	if t.Before(today) {
		if task.Repeat == "" {
			// Если нет правила повторения, берем сегодняшнюю дату
			task.Date = now.Format(DateFormat)
		} else {
			// Иначе используем вычисленную следующую дату
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("правило повторения указано в неправильном формате")
			}
			task.Date = next
		}
	}

	// Если дата равна сегодняшней или больше - оставляем как есть

	return nil
}
