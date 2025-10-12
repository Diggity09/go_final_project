package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
	var install bool

	// Проверяем существование файла
	if _, err := os.Stat(dbFile); err != nil {
		install = true
	}

	// Открываем базу данных
	var err error
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Если файл не существовал, создаем таблицы
	if install {
		if _, err := db.Exec(schema); err != nil {
			return err
		}
	}

	return nil
}
