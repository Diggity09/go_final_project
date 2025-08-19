package api

import (
	"net/http"
	"os"
)

func Init() {
	// API обработчики
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", doneTaskHandler)

	// Обслуживание статических файлов
	setupStaticFiles()
}

func setupStaticFiles() {
	// Проверяем, какая папка существует для статических файлов
	var staticDir string

	if _, err := os.Stat("./web"); err == nil {
		staticDir = "./web"
	} else if _, err := os.Stat("./static"); err == nil {
		staticDir = "./static"
	} else if _, err := os.Stat("./public"); err == nil {
		staticDir = "./public"
	} else {
		// Создаем простую заглушку
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.Write([]byte(`
<!DOCTYPE html>
<html>
<head><title>Планировщик задач</title></head>
<body>
<h1>Планировщик задач запущен</h1>
<p>API доступно по адресам:</p>
<ul>
<li>POST /api/task - добавить задачу</li>
<li>GET /api/task?id=X - получить задачу по ID</li>
<li>PUT /api/task - обновить задачу</li>
<li>DELETE /api/task?id=X - удалить задачу</li>
<li>POST /api/task/done?id=X - завершить задачу</li>
<li>GET /api/tasks - получить список задач</li>
<li>GET /api/nextdate - вычислить следующую дату</li>
</ul>
</body>
</html>`))
			} else {
				http.NotFound(w, r)
			}
		})
		return
	}

	// Обслуживание файлов из найденной папки
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fs)
}

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
