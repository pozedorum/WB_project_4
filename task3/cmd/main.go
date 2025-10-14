package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pozedorum/WB_project_2/task18/config"
	"github.com/pozedorum/WB_project_2/task18/internal/server"
	"github.com/pozedorum/WB_project_2/task18/internal/service"
	"github.com/pozedorum/WB_project_2/task18/internal/storage"
)

func main() {

	cfg := config.Load()
	// Инициализация зависимостей
	repo := storage.NewEventStorage()
	serv := service.NewEventService(repo)
	srv, err := server.NewServer(*cfg, serv)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Server started on :%s", cfg.Port)
	if err := srv.Run(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
