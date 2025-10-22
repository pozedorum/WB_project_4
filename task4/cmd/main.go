package main

import (
	"log"
	"os"
	"strconv"

	"github.com/pozedorum/WB_project_4/task4/internal/server"
)

func main() {
	// Получаем порт из переменной окружения или используем по умолчанию
	port := getPort()

	// Создаем сервер
	serv := server.NewServer()

	log.Printf("🚀 Starting Go Metrics Server on :%s", port)
	log.Printf("📊 Web Dashboard: http://localhost:%s", port)
	log.Printf("🔧 Prometheus Metrics: http://localhost:%s/metrics", port)
	log.Printf("📚 API Base URL: http://localhost:%s/api", port)

	// Запускаем сервер
	addr := ":" + port
	if err := serv.Run(addr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

// getPort возвращает порт из переменной окружения или "8080" по умолчанию
func getPort() string {
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if _, err := strconv.Atoi(port); err == nil {
			return port
		}
		log.Printf("⚠️  Invalid PORT environment variable: %s, using default 8080", port)
	}
	return "8080"
}
