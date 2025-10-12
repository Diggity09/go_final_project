package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
)

func main() {
	// Получаем конфигурацию
	port := 7540
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	dbFile := "scheduler.db"
	if envDBFile := os.Getenv("TODO_DBFILE"); envDBFile != "" {
		dbFile = envDBFile
	}

	// Инициализируем базу данных
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	// Инициализируем API
	api.Init()

	// Обслуживаем статические файлы
	http.Handle("/", http.FileServer(http.Dir("./web/")))

	// Запускаем сервер
	log.Printf("Сервер запущен на порту %d", port)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
