package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pozedorum/WB_project_4/task3/internal/models"
)

func (s *EventServer) handleCreateEvent(c *gin.Context) {
	start := time.Now()
	var (
		req  models.EventCreateRequest
		resp models.EventResponse
	)

	s.logger.Debug("HANDLER_CREATE_EVENT", "Starting create event request")

	// Поддержка JSON и form-data
	if err := c.ShouldBind(&req); err != nil {
		s.logger.Warn("HANDLER_CREATE_EVENT", "Invalid input data",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds())
		c.JSON(http.StatusBadRequest, gin.H{"error": models.Err400InvalidInput.Error()})
		return
	}

	// Валидация
	if req.UserID == 0 {
		s.logger.Warn("HANDLER_CREATE_EVENT", "Missing user_id")
		c.JSON(http.StatusBadRequest, gin.H{"error": models.Err400EmptyUserID.Error()})
		return
	}
	if req.Text == "" {
		s.logger.Warn("HANDLER_CREATE_EVENT", "Missing text")
		c.JSON(http.StatusBadRequest, gin.H{"error": models.Err400EmptyText.Error()})
		return
	}
	if req.Datetime.IsZero() {
		s.logger.Warn("HANDLER_CREATE_EVENT", "Missing datetime")
		c.JSON(http.StatusBadRequest, gin.H{"error": models.Err400EmptyDatetime.Error()})
		return
	}

	s.logger.Info("HANDLER_CREATE_EVENT", "Processing event creation",
		"user_id", req.UserID,
		"title", req.Title,
		"datetime", req.Datetime)

	if resp = s.serv.CreateEvent(req); resp.Error != nil {
		s.logger.Error("HANDLER_CREATE_EVENT", "Service layer error",
			"error", resp.Error,
			"user_id", req.UserID,
			"title", req.Title)

		switch resp.Error {
		case models.Err503AlreadyExists:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": resp.Error.Error()})
		case models.Err503PastDate:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": resp.Error.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		}
		return
	}

	duration := time.Since(start)
	s.logger.Info("HANDLER_CREATE_EVENT", "Event created successfully",
		"user_id", req.UserID,
		"event_id", resp.EventID,
		"duration_ms", duration.Milliseconds())

	c.JSON(http.StatusOK, gin.H{
		"result": fmt.Sprintf("event created successfully: %s", req.Text),
	})
}

func (s *EventServer) handleUpdateEvent(c *gin.Context) {
	start := time.Now()
	var (
		req  models.EventUpdateRequest
		resp models.EventResponse
	)

	s.logger.Debug("HANDLER_UPDATE_EVENT", "Starting update event request")

	if err := c.ShouldBind(&req); err != nil {
		s.logger.Warn("HANDLER_UPDATE_EVENT", "Invalid input data",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds())
		c.JSON(http.StatusBadRequest, gin.H{"error": models.Err400InvalidInput.Error()})
		return
	}

	// Валидация
	if req.EventID == 0 {
		s.logger.Warn("HANDLER_UPDATE_EVENT", "Missing event ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": models.Err400InvalidEventID.Error()})
		return
	}

	s.logger.Info("HANDLER_UPDATE_EVENT", "Processing event update",
		"event_id", req.EventID,
		"title", req.Title)

	if resp = s.serv.UpdateEvent(req); resp.Error != nil {
		s.logger.Error("HANDLER_UPDATE_EVENT", "Service layer error",
			"error", resp.Error,
			"event_id", req.EventID,
			"title", req.Title)

		switch resp.Error {
		case models.Err503NotFound:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": resp.Error.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		}
		return
	}

	duration := time.Since(start)
	s.logger.Info("HANDLER_UPDATE_EVENT", "Event updated successfully",
		"event_id", resp.EventID,
		"user_id", resp.UserID,
		"duration_ms", duration.Milliseconds())

	c.JSON(http.StatusOK, gin.H{
		"result": fmt.Sprintf("event #%d updated successfully", resp.EventID),
	})
}

func (s *EventServer) handleDeleteEvent(c *gin.Context) {
	var req models.EventDeleteRequest
	start := time.Now()

	s.logger.Debug("HANDLER_DELETE_EVENT", "Starting delete event request")

	if err := c.ShouldBind(&req); err != nil {
		s.logger.Warn("HANDLER_DELETE_EVENT", "Invalid input data",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds())
		c.JSON(http.StatusBadRequest, gin.H{"error": models.Err400InvalidInput.Error()})
		return
	}

	s.logger.Info("HANDLER_DELETE_EVENT", "Processing event deletion",
		"event_id", req.EventID)

	if err := s.serv.DeleteEvent(req); err != nil {
		s.logger.Error("HANDLER_DELETE_EVENT", "Service layer error",
			"error", err,
			"event_id", req.EventID)

		switch err {
		case models.Err503NotFound:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		}
		return
	}

	duration := time.Since(start)
	s.logger.Info("HANDLER_DELETE_EVENT", "Event deleted successfully",
		"event_id", req.EventID,
		"duration_ms", duration.Milliseconds())

	c.JSON(http.StatusOK, gin.H{
		"result": fmt.Sprintf("event #%s deleted successfully", req.EventID),
	})
}

func (s *EventServer) handleGetDayEvents(c *gin.Context) {
	start := time.Now()

	s.logger.Debug("HANDLER_GET_DAY_EVENTS", "Starting get day events request")

	userID, date, err := s.parseQueryParams(c)
	if err != nil {
		s.logger.Warn("HANDLER_GET_DAY_EVENTS", "Invalid query parameters",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req := models.EventsGetRequest{UserID: userID, Date: date}
	s.logger.Info("HANDLER_GET_DAY_EVENTS", "Processing day events request",
		"user_id", userID,
		"date", date)

	events, err := s.serv.GetDayEvents(req)
	if err != nil {
		s.logger.Error("HANDLER_GET_DAY_EVENTS", "Service layer error",
			"error", err,
			"user_id", userID,
			"date", date)
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		return
	}

	duration := time.Since(start)
	s.logger.Info("HANDLER_GET_DAY_EVENTS", "Day events retrieved successfully",
		"user_id", userID,
		"date", date,
		"events_count", len(events),
		"duration_ms", duration.Milliseconds())

	c.JSON(http.StatusOK, gin.H{"result": events})
}

func (s *EventServer) handleGetWeekEvents(c *gin.Context) {
	start := time.Now()

	s.logger.Debug("HANDLER_GET_WEEK_EVENTS", "Starting get week events request")

	userID, date, err := s.parseQueryParams(c)
	if err != nil {
		s.logger.Warn("HANDLER_GET_WEEK_EVENTS", "Invalid query parameters",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req := models.EventsGetRequest{UserID: userID, Date: date}
	s.logger.Info("HANDLER_GET_WEEK_EVENTS", "Processing week events request",
		"user_id", req.UserID,
		"date", req.Date)

	events, err := s.serv.GetWeekEvents(req)
	if err != nil {
		s.logger.Error("HANDLER_GET_WEEK_EVENTS", "Service layer error",
			"error", err,
			"user_id", req.UserID,
			"date", req.Date)
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		return
	}

	duration := time.Since(start)
	s.logger.Info("HANDLER_GET_WEEK_EVENTS", "Week events retrieved successfully",
		"user_id", req.UserID,
		"date", req.Date,
		"events_count", len(events),
		"duration_ms", duration.Milliseconds())

	c.JSON(http.StatusOK, gin.H{"result": events})
}

func (s *EventServer) handleGetMonthEvents(c *gin.Context) {
	start := time.Now()

	s.logger.Debug("HANDLER_GET_MONTH_EVENTS", "Starting get month events request")

	userID, date, err := s.parseQueryParams(c)
	if err != nil {
		s.logger.Warn("HANDLER_GET_MONTH_EVENTS", "Invalid query parameters",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req := models.EventsGetRequest{UserID: userID, Date: date}
	s.logger.Info("HANDLER_GET_MONTH_EVENTS", "Processing month events request",
		"user_id", req.UserID,
		"date", req.Date)

	events, err := s.serv.GetMonthEvents(req)
	if err != nil {
		s.logger.Error("HANDLER_GET_MONTH_EVENTS", "Service layer error",
			"error", err,
			"user_id", req.UserID,
			"date", req.Date)
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		return
	}

	duration := time.Since(start)
	s.logger.Info("HANDLER_GET_MONTH_EVENTS", "Month events retrieved successfully",
		"user_id", req.UserID,
		"date", req.Date,
		"events_count", len(events),
		"duration_ms", duration.Milliseconds())

	c.JSON(http.StatusOK, gin.H{"result": events})
}

// Вспомогательные методы
func (s *EventServer) parseQueryParams(c *gin.Context) (int, time.Time, error) {
	userID, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("invalid userID: %v", err)
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		return 0, time.Time{}, models.ErrEmptyDatetime
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("invalid date format: %w", err)
	}

	return userID, date, nil
}
