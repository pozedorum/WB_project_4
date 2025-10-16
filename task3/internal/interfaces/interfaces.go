package interfaces

import (
	"context"
	"time"

	"github.com/pozedorum/WB_project_4/task3/internal/models"
)

type Server interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type Service interface {
	CreateEvent(event models.Event) error
	UpdateEvent(event models.Event) error
	DeleteEvent(event models.Event) error
	GetDayEvents(userID string, date time.Time) ([]models.Event, error)
	GetWeekEvents(userID string, startWeek time.Time) ([]models.Event, error)
	GetMonthEvents(userID string, startMonth time.Time) ([]models.Event, error)
}

type Repository interface {
	CreateEvent(event models.Event) error
	UpdateEvent(event models.Event) error
	DeleteEvent(event models.Event) error
	GetByDateRange(start, end time.Time) ([]models.Event, error)
	GetEventByID(id int) (*models.Event, error)
	Close() error
}
