// Package interfaces содержит все интерфейсы используемые в DI
package interfaces

import (
	"context"
	"time"

	"github.com/pozedorum/WB_project_4/task3/internal/models"
)

// Server интерфейс HTTP сервера
type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
}

// Service интерфейс бизнес-логики событий
type Service interface {
	CreateEvent(req models.EventCreateRequest) models.EventResponse
	UpdateEvent(req models.EventUpdateRequest) models.EventResponse
	DeleteEvent(req models.EventDeleteRequest) error
	GetDayEvents(req models.EventsGetRequest) ([]models.Event, error)
	GetWeekEvents(req models.EventsGetRequest) ([]models.Event, error)
	GetMonthEvents(req models.EventsGetRequest) ([]models.Event, error)
}

// Repository интерфейс работы с данными
type Repository interface {
	CreateEvent(event models.Event) error
	UpdateEvent(event models.Event) error
	DeleteEvent(event models.Event) error
	GetByDateRange(start, end time.Time) ([]models.Event, error)
	GetEventByID(id int) (*models.Event, error)
	Close() error

	ArchiveOldEvents(threshold time.Time) (int, error)
	GetArchivedEvents() ([]models.Event, error)
	CleanupArchivedEvents(beforeTime time.Time) (int, error)
}

// Reminder определяет интерфейс сервиса напоминаний
type Reminder interface {
	ScheduleReminder(event models.Event) error
	UpdateReminder(event models.Event) error
	CancelReminder(eventID int) error
	StartWorker(ctx context.Context) error
	Shutdown()
}

type Archiver interface {
	Start(ctx context.Context)
}

// TelegramClient определяет интерфейс для Telegram клиента
type TelegramClient interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
	GetBotInfo() string
	Close()
}

// QueueClient определяет интерфейс для работы с очередью
type QueueClient interface {
	PublishReminder(reminder models.ReminderMessage) error
	StartConsuming(ctx context.Context, handler func(models.ReminderMessage) error) error
	Close()
}

// Logger интерфейс логгера
type Logger interface {
	Debug(operation, message string, keyvals ...interface{})
	Info(operation, message string, keyvals ...interface{})
	Warn(operation, message string, keyvals ...interface{})
	Error(operation, message string, keyvals ...interface{})
	Shutdown()
}
