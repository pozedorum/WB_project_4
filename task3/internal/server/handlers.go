package server

import (
	"fmt"
	"net/http"
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
	if req.UserName == "" {
		s.logger.Warn("HANDLER_CREATE_EVENT", "Missing username")
		c.JSON(http.StatusBadRequest, gin.H{"error": models.Err400EmptyUserName.Error()})
		return
	}
	// пустой текст допустим
	// if req.Text == "" {
	// 	s.logger.Warn("HANDLER_CREATE_EVENT", "Missing text")
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": models.Err400EmptyText.Error()})
	// 	return
	// }
	if req.Datetime.IsZero() {
		s.logger.Warn("HANDLER_CREATE_EVENT", "Missing datetime")
		c.JSON(http.StatusBadRequest, gin.H{"error": models.Err400EmptyDatetime.Error()})
		return
	}

	s.logger.Info("HANDLER_CREATE_EVENT", "Processing event creation",
		"username", req.UserName,
		"title", req.Title,
		"datetime", req.Datetime)

	if resp = s.serv.CreateEvent(req); resp.Error != nil {
		s.logger.Error("HANDLER_CREATE_EVENT", "Service layer error",
			"error", resp.Error,
			"username", req.UserName,
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
		"username", req.UserName,
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
		"username", resp.UserName,
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
		"result": fmt.Sprintf("event #%d deleted successfully", req.EventID),
	})
}

func (s *EventServer) handleGetDayEvents(c *gin.Context) {
	start := time.Now()

	s.logger.Debug("HANDLER_GET_DAY_EVENTS", "Starting get day events request")

	userName, date, err := s.parseQueryParams(c)
	if err != nil {
		s.logger.Warn("HANDLER_GET_DAY_EVENTS", "Invalid query parameters",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req := models.EventsGetRequest{UserName: userName, Date: date}
	s.logger.Info("HANDLER_GET_DAY_EVENTS", "Processing day events request",
		"username", userName,
		"date", date)

	events, err := s.serv.GetDayEvents(req)
	if err != nil {
		s.logger.Error("HANDLER_GET_DAY_EVENTS", "Service layer error",
			"error", err,
			"username", userName,
			"date", date)
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		return
	}

	duration := time.Since(start)
	s.logger.Info("HANDLER_GET_DAY_EVENTS", "Day events retrieved successfully",
		"username", userName,
		"date", date,
		"events_count", len(events),
		"duration_ms", duration.Milliseconds())

	c.JSON(http.StatusOK, gin.H{"result": events})
}

func (s *EventServer) handleGetWeekEvents(c *gin.Context) {
	start := time.Now()

	s.logger.Debug("HANDLER_GET_WEEK_EVENTS", "Starting get week events request")

	userName, date, err := s.parseQueryParams(c)
	if err != nil {
		s.logger.Warn("HANDLER_GET_WEEK_EVENTS", "Invalid query parameters",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req := models.EventsGetRequest{UserName: userName, Date: date}
	s.logger.Info("HANDLER_GET_WEEK_EVENTS", "Processing week events request",
		"username", req.UserName,
		"date", req.Date)

	events, err := s.serv.GetWeekEvents(req)
	if err != nil {
		s.logger.Error("HANDLER_GET_WEEK_EVENTS", "Service layer error",
			"error", err,
			"username", req.UserName,
			"date", req.Date)
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		return
	}

	duration := time.Since(start)
	s.logger.Info("HANDLER_GET_WEEK_EVENTS", "Week events retrieved successfully",
		"username", req.UserName,
		"date", req.Date,
		"events_count", len(events),
		"duration_ms", duration.Milliseconds())

	c.JSON(http.StatusOK, gin.H{"result": events})
}

func (s *EventServer) handleGetMonthEvents(c *gin.Context) {
	start := time.Now()

	s.logger.Debug("HANDLER_GET_MONTH_EVENTS", "Starting get month events request")

	userName, date, err := s.parseQueryParams(c)
	if err != nil {
		s.logger.Warn("HANDLER_GET_MONTH_EVENTS", "Invalid query parameters",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req := models.EventsGetRequest{UserName: userName, Date: date}
	s.logger.Info("HANDLER_GET_MONTH_EVENTS", "Processing month events request",
		"username", req.UserName,
		"date", req.Date)

	events, err := s.serv.GetMonthEvents(req)
	if err != nil {
		s.logger.Error("HANDLER_GET_MONTH_EVENTS", "Service layer error",
			"error", err,
			"username", req.UserName,
			"date", req.Date)
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		return
	}

	duration := time.Since(start)
	s.logger.Info("HANDLER_GET_MONTH_EVENTS", "Month events retrieved successfully",
		"username", req.UserName,
		"date", req.Date,
		"events_count", len(events),
		"duration_ms", duration.Milliseconds())

	c.JSON(http.StatusOK, gin.H{"result": events})
}

// Вспомогательные методы
func (s *EventServer) parseQueryParams(c *gin.Context) (string, time.Time, error) {
	userName := c.Query("username")
	// if err != nil {
	// 	return 0, time.Time{}, fmt.Errorf("invalid userName: %v", err)
	// }

	dateStr := c.Query("date")
	if dateStr == "" {
		return "", time.Time{}, models.ErrEmptyDatetime
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid date format: %w", err)
	}

	return userName, date, nil
}
