package api

// Task представляет структуру задачи для API
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

// AddTaskResponse представляет ответ при добавлении задачи
type AddTaskResponse struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// UpdateTaskResponse представляет ответ при обновлении задачи
type UpdateTaskResponse struct {
	Error string `json:"error,omitempty"`
}

// TasksResponse представляет ответ со списком задач
type TasksResponse struct {
	Tasks []*Task `json:"tasks"`
}

// NextDateRequest представляет запрос для вычисления следующей даты
type NextDateRequest struct {
	Now    string `json:"now"`
	Date   string `json:"date"`
	Repeat string `json:"repeat"`
}

// NextDateResponse представляет ответ с следующей датой
type NextDateResponse struct {
	Date  string `json:"date,omitempty"`
	Error string `json:"error,omitempty"`
}
