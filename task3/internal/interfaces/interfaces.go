package interfaces

import (
	"context"
	"time"

	"github.com/pozedorum/WB_project_4/task3/internal/models"
)

type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
}

type Service interface {
	CreateEvent(req models.EventCreateRequest) models.EventResponse
	UpdateEvent(req models.EventUpdateRequest) models.EventResponse
	DeleteEvent(req models.EventDeleteRequest) error
	GetDayEvents(req models.EventsGetRequest) ([]models.Event, error)
	GetWeekEvents(req models.EventsGetRequest) ([]models.Event, error)
	GetMonthEvents(req models.EventsGetRequest) ([]models.Event, error)
}

type Repository interface {
	CreateEvent(event models.Event) error
	UpdateEvent(event models.Event) error
	DeleteEvent(event models.Event) error
	GetByDateRange(start, end time.Time) ([]models.Event, error)
	GetEventByID(id int) (*models.Event, error)
	Close() error
}

type Logger interface {
	Debug(operation, message string, keyvals ...interface{})
	Info(operation, message string, keyvals ...interface{})
	Warn(operation, message string, keyvals ...interface{})
	Error(operation, message string, keyvals ...interface{})
	Shutdown()
}
