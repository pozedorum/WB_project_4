package di

import (
	"context"
	"fmt"
	"time"

	"github.com/pozedorum/WB_project_4/task3/internal/interfaces"
	"github.com/pozedorum/WB_project_4/task3/internal/repository"
	"github.com/pozedorum/WB_project_4/task3/internal/server"
	"github.com/pozedorum/WB_project_4/task3/internal/service"
	"github.com/pozedorum/WB_project_4/task3/pkg/config"
	"github.com/pozedorum/WB_project_4/task3/pkg/logger"
)

type Container struct {
	repo    interfaces.Repository
	service interfaces.Service
	server  interfaces.Server
	logger  interfaces.Logger
}

func NewContainer(cfg *config.Config) (*Container, error) {
	// Инициализируем логгер
	logger, err := logger.NewLogger("event-service", "")
	//	logger, err := logger.NewLogger("event-service", "./logs/app.log")
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	logger.Info("CONTAINER_INIT", "Starting application container initialization")

	// Передаем логгер в репозиторий
	repo, err := repository.NewEventRepository(cfg.Database.GetDSN(), logger)
	if err != nil {
		logger.Error("CONTAINER_INIT", "Failed to create repository", "error", err)
		return nil, err
	}
	logger.Info("CONTAINER_INIT", "Repository initialized successfully")

	// Передаем логгер в сервис
	service := service.NewEventService(repo, logger)
	logger.Info("CONTAINER_INIT", "Service initialized successfully")

	server := server.NewEventServer(cfg.Server.Port, service, logger)
	logger.Info("CONTAINER_INIT", "Server initialized successfully")

	return &Container{
		repo:    repo,
		service: service,
		server:  server,
		logger:  logger,
	}, nil
}

func (c *Container) Start() error {
	return c.server.Start()
}

func (c *Container) Shutdown() error {
	var errors []error
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := c.server.Shutdown(ctx)
	if err != nil {
		errors = append(errors, err)
	}
	err = c.repo.Close()
	if err != nil {
		errors = append(errors, err)
	}
	if len(errors) > 0 {
		return fmt.Errorf("shutdown completed with errors: %v", errors)
	}
	return nil
}
