package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

// Task представляет структуру задачи
type Task struct {
	ID      int64  `json:"id" db:"id"`
	Date    string `json:"date" db:"date"`
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"`
}

// GetTask возвращает задачу по ID
func GetTask(id string) (*Task, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	// Преобразуем строковый ID в int64
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("некорректный идентификатор")
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	var task Task

	err = db.QueryRow(query, idInt).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, err
	}

	return &task, nil
}

// UpdateTask обновляет существующую задачу
func UpdateTask(task *Task) error {
	if db == nil {
		return sql.ErrConnDone
	}

	// Преобразуем строковый ID в int64 если нужно
	var idInt int64
	var err error

	if task.ID == 0 {
		return fmt.Errorf("некорректный идентификатор")
	}
	idInt = task.ID

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, idInt)
	if err != nil {
		return err
	}

	// Проверяем количество обновленных записей
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}

	return nil
}

// UpdateTaskDate обновляет только дату задачи
func UpdateTaskDate(newDate, id string) error {
	if db == nil {
		return sql.ErrConnDone
	}

	// Преобразуем строковый ID в int64
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("некорректный идентификатор")
	}

	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := db.Exec(query, newDate, idInt)
	if err != nil {
		return err
	}

	// Проверяем количество обновленных записей
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

// DeleteTask удаляет задачу по ID
func DeleteTask(id string) error {
	if db == nil {
		return sql.ErrConnDone
	}

	// Преобразуем строковый ID в int64
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("некорректный идентификатор")
	}

	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.Exec(query, idInt)
	if err != nil {
		return err
	}

	// Проверяем количество удаленных записей
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

// Tasks возвращает список задач с ограничением по количеству
func Tasks(limit int) ([]*Task, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Возвращаем пустой слайс, если задач нет
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// GetTasks возвращает все задачи, отсортированные по дате
func GetTasks() ([]Task, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

// AddTask добавляет новую задачу в базу данных
func AddTask(task *Task) (int64, error) {
	if db == nil {
		return 0, sql.ErrConnDone
	}

	// Если дата пустая, используем сегодняшнюю дату
	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// TaskExists проверяет, существует ли задача с указанным ID
func TaskExists(id int64) (bool, error) {
	if db == nil {
		return false, sql.ErrConnDone
	}

	query := `SELECT COUNT(*) FROM scheduler WHERE id = ?`
	var count int
	err := db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
