package server

import (
	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
	"log"
	"net/http"
)

func Start() error {
	// Инициализация базы данных
	if err := db.Init(); err != nil {
		log.Fatal("Ошибка инициализации базы данных:", err)
		return err
	}

	// Инициализация API обработчиков
	api.Init()

	// Запуск сервера
	log.Println("Сервер запущен на порту 7540")
	return http.ListenAndServe(":7540", nil)
}
