package reminder

import (
	"context"
	"fmt"
	"time"

	"github.com/pozedorum/WB_project_4/task3/internal/interfaces"
	"github.com/pozedorum/WB_project_4/task3/internal/models"
)

type reminderService struct {
	queueClient    interfaces.QueueClient
	telegramClient interfaces.TelegramClient
	logger         interfaces.Logger
	shutdownChan   chan struct{}
}

type Config struct {
	RabbitMQURL   string
	QueueName     string
	ExchangeName  string
	TelegramToken string
}

func NewService(cfg Config, logger interfaces.Logger) (interfaces.Reminder, error) {
	// Инициализация RabbitMQ клиента
	queueClient, err := NewRabbitMQClient(RabbitMQConfig{
		URL:       cfg.RabbitMQURL,
		QueueName: cfg.QueueName,
		Exchange:  cfg.ExchangeName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ client: %w", err)
	}

	// Инициализация Telegram клиента
	telegramClient, err := NewTelegramClient(cfg.TelegramToken)
	if err != nil {
		queueClient.Close()
		return nil, fmt.Errorf("failed to create Telegram client: %w", err)
	}

	service := &reminderService{
		queueClient:    queueClient,
		telegramClient: telegramClient,
		logger:         logger,
		shutdownChan:   make(chan struct{}),
	}

	logger.Info("REMINDER_SERVICE", "Service initialized successfully",
		"queue", cfg.QueueName, "exchange", cfg.ExchangeName, "telegram_bot", telegramClient.GetBotInfo())

	return service, nil
}

func (s *reminderService) ScheduleReminder(event models.Event) error {
	if event.RemindBefore <= 0 || event.TelegramID == 0 {
		return nil // Напоминание не требуется
	}

	remindTime := event.Datetime.Add(-time.Duration(event.RemindBefore) * time.Minute)

	message := models.ReminderMessage{
		EventID:    event.ID,
		UserToken:  event.UserToken,
		TelegramID: event.TelegramID,
		Title:      event.Title,
		Message:    fmt.Sprintf("Напоминание: %s\n%s", event.Title, event.Text),
		NotifyTime: remindTime,
	}

	if err := s.queueClient.PublishReminder(message); err != nil {
		return fmt.Errorf("failed to publish reminder: %w", err)
	}

	s.logger.Info("REMINDER_SCHEDULE", "Reminder scheduled successfully",
		"event_id", event.ID, "notify_time", remindTime, "telegram_id", event.TelegramID)

	return nil
}

func (s *reminderService) UpdateReminder(event models.Event) error {
	// Для простоты пересоздаем напоминание
	return s.ScheduleReminder(event)
}

func (s *reminderService) CancelReminder(eventID int) error {
	s.logger.Info("REMINDER_CANCEL", "Reminder cancellation requested", "event_id", eventID)
	// В production системе здесь была бы логика удаления из очереди
	return nil
}

func (s *reminderService) StartWorker(ctx context.Context) error {
	handler := func(reminder models.ReminderMessage) error {
		return s.processReminder(ctx, reminder)
	}

	go func() {
		s.logger.Info("REMINDER_WORKER", "Starting reminder worker")

		if err := s.queueClient.StartConsuming(ctx, handler); err != nil {
			s.logger.Error("REMINDER_WORKER", "Worker stopped with error", "error", err)
		}

		s.logger.Info("REMINDER_WORKER", "Reminder worker stopped")
	}()

	return nil
}

func (s *reminderService) processReminder(ctx context.Context, reminder models.ReminderMessage) error {
	// Создаем контекст с таймаутом для Telegram
	telegramCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Отправляем уведомление
	if err := s.telegramClient.SendMessage(telegramCtx, reminder.TelegramID, reminder.Message); err != nil {
		return fmt.Errorf("failed to send Telegram notification: %w", err)
	}

	s.logger.Info("REMINDER_PROCESSED", "Reminder sent successfully",
		"event_id", reminder.EventID, "telegram_id", reminder.TelegramID)

	return nil
}

func (s *reminderService) Shutdown() {
	s.logger.Info("REMINDER_SHUTDOWN", "Shutting down reminder service")
	close(s.shutdownChan)
	s.queueClient.Close()
	s.telegramClient.Close()
	s.logger.Info("REMINDER_SHUTDOWN", "Reminder service shutdown completed")
}
