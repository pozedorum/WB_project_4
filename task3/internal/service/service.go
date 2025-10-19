package service

import (
	"time"

	"github.com/pozedorum/WB_project_4/task3/internal/interfaces"
	"github.com/pozedorum/WB_project_4/task3/internal/models"
)

type EventService struct {
	repo     interfaces.Repository
	reminder interfaces.Reminder
	logger   interfaces.Logger
}

func NewEventService(repo interfaces.Repository, reminder interfaces.Reminder, logger interfaces.Logger) *EventService {
	return &EventService{
		repo:     repo,
		reminder: reminder,
		logger:   logger,
	}
}

func (s *EventService) CreateEvent(req models.EventCreateRequest) models.EventResponse {
	start := time.Now()
	response := models.EventResponse{}

	s.logger.Debug("SERVICE_CREATE_EVENT", "Starting event creation",
		"usertoken", req.UserToken,
		"title", req.Title,
		"datetime", req.Datetime)

	// Бизнес-логика: нельзя создавать события в прошлом
	if req.Datetime.Before(time.Now()) {
		s.logger.Warn("SERVICE_CREATE_EVENT", "Attempt to create event in the past",
			"usertoken", req.UserToken,
			"title", req.Title,
			"datetime", req.Datetime,
			"duration_ms", time.Since(start).Milliseconds())
		response.Error = models.Err503PastDate
		return response
	}

	event := models.Event{
		UserToken:    req.UserToken,
		Title:        req.Title,
		Text:         req.Text,
		Datetime:     req.Datetime,
		RemindBefore: int(req.RemindBefore / time.Second),
		IsArchived:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := s.repo.CreateEvent(event)
	if err != nil {
		s.logger.Error("SERVICE_CREATE_EVENT", "Failed to create event",
			"error", err,
			"usertoken", req.UserToken,
			"title", req.Title,
			"duration_ms", time.Since(start).Milliseconds())
		response.Error = err
		return response
	}

	if event.RemindBefore > 0 && req.TelegramID != 0 {
		if err := s.reminder.ScheduleReminder(event); err != nil {
			s.logger.Warn("SERVICE_CREATE_EVENT", "Failed to schedule reminder",
				"event_id", event.ID, "error", err)
		}
	}

	s.logger.Info("SERVICE_CREATE_EVENT", "Event created successfully",
		"usertoken", req.UserToken,
		"event_id", event.ID,
		"duration_ms", time.Since(start).Milliseconds())

	response.UserToken = req.UserToken
	response.EventID = event.ID
	response.Title = req.Title
	response.EventDatetime = req.Datetime
	return response
}

func (s *EventService) UpdateEvent(req models.EventUpdateRequest) models.EventResponse {
	start := time.Now()
	response := models.EventResponse{}

	s.logger.Debug("SERVICE_UPDATE_EVENT", "Starting event update",
		"event_id", req.EventID,
		"title", req.Title)

	// Проверяем существование события
	existing, err := s.repo.GetEventByID(req.EventID)
	if err != nil {
		s.logger.Error("SERVICE_UPDATE_EVENT", "Event not found for update",
			"error", err,
			"event_id", req.EventID,
			"duration_ms", time.Since(start).Milliseconds())
		response.Error = err
		return response
	}

	if req.Datetime.Before(time.Now()) {
		s.logger.Warn("SERVICE_UPDATE_EVENT", "Attempt to move event to the past",
			"event_id", req.EventID,
			"datetime", req.Datetime,
			"duration_ms", time.Since(start).Milliseconds())
		response.Error = models.Err503PastDate
		return response
	}

	event := models.Event{
		ID:           req.EventID,
		UserToken:    existing.UserToken, // UserToken не меняется при обновлении
		Title:        req.Title,
		Text:         req.Text,
		Datetime:     req.Datetime,
		RemindBefore: int(req.RemindBefore / time.Second),
		IsArchived:   existing.IsArchived,
		CreatedAt:    existing.CreatedAt,
		UpdatedAt:    time.Now(),
	}

	err = s.repo.UpdateEvent(event)
	if err != nil {
		s.logger.Error("SERVICE_UPDATE_EVENT", "Failed to update event",
			"error", err,
			"event_id", req.EventID,
			"duration_ms", time.Since(start).Milliseconds())
		response.Error = err
		return response
	}

	if event.RemindBefore > 0 && existing.TelegramID != 0 {
		if err := s.reminder.UpdateReminder(event); err != nil {
			s.logger.Warn("SERVICE_UPDATE_EVENT", "Failed to update reminder",
				"event_id", event.ID, "error", err)
		}
	}
	s.logger.Info("SERVICE_UPDATE_EVENT", "Event updated successfully",
		"event_id", req.EventID,
		"usertoken", existing.UserToken,
		"duration_ms", time.Since(start).Milliseconds())

	response.UserToken = existing.UserToken
	response.EventID = req.EventID
	response.Title = req.Title
	response.EventDatetime = req.Datetime
	return response
}

func (s *EventService) DeleteEvent(req models.EventDeleteRequest) error {
	start := time.Now()

	s.logger.Debug("SERVICE_DELETE_EVENT", "Starting event deletion",
		"event_id", req.EventID)

	event := models.Event{
		ID: req.EventID,
	}

	err := s.repo.DeleteEvent(event)
	if err != nil {
		s.logger.Error("SERVICE_DELETE_EVENT", "Failed to delete event",
			"error", err,
			"event_id", req.EventID,
			"duration_ms", time.Since(start).Milliseconds())
		return err
	}

	s.logger.Info("SERVICE_DELETE_EVENT", "Event deleted successfully",
		"event_id", req.EventID,
		"duration_ms", time.Since(start).Milliseconds())

	return nil
}

func (s *EventService) GetDayEvents(req models.EventsGetRequest) ([]models.Event, error) {
	start := time.Now()

	s.logger.Debug("SERVICE_GET_DAY_EVENTS", "Getting day events",
		"usertoken", req.UserToken,
		"date", req.Date)

	startRange := time.Date(req.Date.Year(), req.Date.Month(), req.Date.Day(), 0, 0, 0, 0, req.Date.Location())
	end := startRange.Add(24 * time.Hour)

	events, err := s.repo.GetByDateRange(startRange, end)
	if err != nil {
		s.logger.Error("SERVICE_GET_DAY_EVENTS", "Failed to get day events",
			"error", err,
			"usertoken", req.UserToken,
			"date", req.Date,
			"duration_ms", time.Since(start).Milliseconds())
		return nil, err
	}

	filteredEvents := s.filterEventsByUserToken(events, req.UserToken)

	s.logger.Info("SERVICE_GET_DAY_EVENTS", "Day events retrieved successfully",
		"usertoken", req.UserToken,
		"date", req.Date,
		"total_events", len(events),
		"filtered_events", len(filteredEvents),
		"duration_ms", time.Since(start).Milliseconds())

	return filteredEvents, nil
}

func (s *EventService) GetWeekEvents(req models.EventsGetRequest) ([]models.Event, error) {
	start := time.Now()

	s.logger.Debug("SERVICE_GET_WEEK_EVENTS", "Getting week events",
		"usertoken", req.UserToken,
		"date", req.Date)

	startRange := time.Date(req.Date.Year(), req.Date.Month(), req.Date.Day(), 0, 0, 0, 0, req.Date.Location())
	end := startRange.Add(7 * 24 * time.Hour)

	events, err := s.repo.GetByDateRange(startRange, end)
	if err != nil {
		s.logger.Error("SERVICE_GET_WEEK_EVENTS", "Failed to get week events",
			"error", err,
			"usertoken", req.UserToken,
			"date", req.Date,
			"duration_ms", time.Since(start).Milliseconds())
		return nil, err
	}

	filteredEvents := s.filterEventsByUserToken(events, req.UserToken)

	s.logger.Info("SERVICE_GET_WEEK_EVENTS", "Week events retrieved successfully",
		"usertoken", req.UserToken,
		"date", req.Date,
		"total_events", len(events),
		"filtered_events", len(filteredEvents),
		"duration_ms", time.Since(start).Milliseconds())

	return filteredEvents, nil
}

func (s *EventService) GetMonthEvents(req models.EventsGetRequest) ([]models.Event, error) {
	start := time.Now()

	s.logger.Debug("SERVICE_GET_MONTH_EVENTS", "Getting month events",
		"usertoken", req.UserToken,
		"date", req.Date)

	startRange := time.Date(req.Date.Year(), req.Date.Month(), 1, 0, 0, 0, 0, req.Date.Location())
	end := startRange.AddDate(0, 1, 0) // Следующий месяц

	events, err := s.repo.GetByDateRange(startRange, end)
	if err != nil {
		s.logger.Error("SERVICE_GET_MONTH_EVENTS", "Failed to get month events",
			"error", err,
			"usertoken", req.UserToken,
			"date", req.Date,
			"duration_ms", time.Since(start).Milliseconds())
		return nil, err
	}

	filteredEvents := s.filterEventsByUserToken(events, req.UserToken)

	s.logger.Info("SERVICE_GET_MONTH_EVENTS", "Month events retrieved successfully",
		"usertoken", req.UserToken,
		"date", req.Date,
		"total_events", len(events),
		"filtered_events", len(filteredEvents),
		"duration_ms", time.Since(start).Milliseconds())

	return filteredEvents, nil
}

func (s *EventService) filterEventsByUserToken(events []models.Event, userToken string) []models.Event {
	var filtered []models.Event
	for _, event := range events {
		if event.UserToken == userToken {
			filtered = append(filtered, event)
		}
	}
	return filtered
}
