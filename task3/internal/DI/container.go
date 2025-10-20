// Package di содержит в себе инициализацию всех зависимостей и слоёв согласно Dependency Injection
package di

import (
	"context"
	"fmt"
	"time"

	"github.com/pozedorum/WB_project_4/task3/internal/archiver"
	"github.com/pozedorum/WB_project_4/task3/internal/interfaces"
	"github.com/pozedorum/WB_project_4/task3/internal/reminder"
	"github.com/pozedorum/WB_project_4/task3/internal/repository"
	"github.com/pozedorum/WB_project_4/task3/internal/server"
	"github.com/pozedorum/WB_project_4/task3/internal/service"
	"github.com/pozedorum/WB_project_4/task3/pkg/config"
	"github.com/pozedorum/WB_project_4/task3/pkg/logger"
)

type Container struct {
	repo     interfaces.Repository
	service  interfaces.Service
	server   interfaces.Server
	reminder interfaces.Reminder
	archiver interfaces.Archiver
	logger   interfaces.Logger
}

func NewContainer(cfg *config.Config) (*Container, error) {
	// Инициализируем логгер
	// logger, err := logger.NewLogger("event-service", "")
	logger, err := logger.NewLogger("event-service", "./logs/app.log")
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	logger.Info("CONTAINER_INIT", "Starting application container initialization")

	// Репозиторий
	repo, err := repository.NewEventRepository(cfg.Database.GetDSN(), logger)
	if err != nil {
		logger.Error("CONTAINER_INIT", "Failed to create repository", "error", err)
		return nil, err
	}
	logger.Info("CONTAINER_INIT", "Repository initialized successfully")

	// Reminder service
	reminderService, err := reminder.NewService(
		cfg.RabbitMQ.URL,
		cfg.RabbitMQ.QueueName,
		cfg.Telegram.Token, logger)
	if err != nil {
		logger.Error("CONTAINER_INIT", "Failed to initialize reminder service", "error", err)
		return nil, err
	}
	logger.Info("CONTAINER_INIT", "Reminder service initialized successfully")

	archiver := archiver.NewArchiveCleaner(repo, logger, 10*time.Minute)
	// Business service
	service := service.NewEventService(repo, reminderService, logger)
	logger.Info("CONTAINER_INIT", "Service initialized successfully")

	// HTTP server
	server := server.NewEventServer(cfg.Server.Port, service, logger)
	logger.Info("CONTAINER_INIT", "Server initialized successfully")

	return &Container{
		repo:     repo,
		service:  service,
		server:   server,
		reminder: reminderService,
		archiver: archiver,
		logger:   logger,
	}, nil
}

func (c *Container) Start() error {
	// Даем время RabbitMQ и PostgreSQL подняться
	time.Sleep(5 * time.Second)
	// ЗАПУСКАЕМ reminder worker перед сервером
	ctx := context.Background()
	if err := c.reminder.StartWorker(ctx); err != nil {
		c.logger.Error("CONTAINER_START", "Failed to start reminder worker", "error", err)
		return err
	}
	c.archiver.Start(ctx)
	c.logger.Info("CONTAINER_START", "Reminder worker started successfully")

	return c.server.Start()
}

func (c *Container) Shutdown() error {
	var errors []error
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := c.server.Shutdown(ctx); err != nil {
		errors = append(errors, fmt.Errorf("server shutdown: %w", err))
	}

	// Shutdown reminder
	c.reminder.Shutdown()

	// Shutdown repository
	if err := c.repo.Close(); err != nil {
		errors = append(errors, fmt.Errorf("repository close: %w", err))
	}

	// Shutdown logger
	c.logger.Shutdown()

	if len(errors) > 0 {
		return fmt.Errorf("shutdown completed with errors: %v", errors)
	}

	c.logger.Info("CONTAINER_SHUTDOWN", "Container shutdown completed successfully")
	return nil
}
