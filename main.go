package main

import (
	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Обработка сигналов для корректного завершения
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Получен сигнал завершения, закрываем приложение...")

		// Закрываем соединение с базой данных
		if err := db.Close(); err != nil {
			log.Printf("Ошибка при закрытии базы данных: %v", err)
		}

		log.Println("Приложение завершено")
		os.Exit(0)
	}()

	// Запуск сервера
	log.Println("Запуск сервера...")
	if err := server.Start(); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
