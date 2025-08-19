package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

// Init инициализирует подключение к базе данных
func Init() error {
	// Путь к базе данных по умолчанию
	dbPath := "scheduler.db"

	// Проверяем, существует ли файл базы данных
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Printf("База данных %s не найдена, будет создана новая", dbPath)
	}

	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	// Проверяем соединение
	if err = db.Ping(); err != nil {
		return err
	}

	// Создаем таблицу, если она не существует
	if err = createTable(); err != nil {
		return err
	}

	log.Println("База данных успешно инициализирована")
	return nil
}

// Close закрывает соединение с базой данных
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// createTable создает таблицу scheduler, если она не существует
func createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat TEXT
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Ошибка создания таблицы: %v", err)
		return err
	}

	log.Println("Таблица scheduler готова к использованию")
	return nil
}

// GetDB возвращает экземпляр базы данных для использования в других пакетах
func GetDB() *sql.DB {
	return db
}
