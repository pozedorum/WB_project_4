package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pozedorum/WB_project_2/task18/config"
	"github.com/pozedorum/WB_project_2/task18/internal/apperrors"
	"github.com/pozedorum/WB_project_2/task18/internal/models"
	"github.com/pozedorum/WB_project_2/task18/pkg/logger"
)

type EventService interface {
	CreateEvent(event models.Event) error
	UpdateEvent(event models.Event) error
	DeleteEvent(event models.Event) error
	GetDayEvents(userID string, date time.Time) ([]models.Event, error)
	GetWeekEvents(userID string, startWeek time.Time) ([]models.Event, error)
	GetMonthEvents(userID string, startMonth time.Time) ([]models.Event, error)
}

type Server struct {
	httpServer *http.Server
	serv       EventService
	logger     *logger.Logger
}

func NewServer(conf config.Config, service EventService) (*Server, error) {
	lg, err := logger.New(conf.FilePath, conf.ConsoleLog)
	if err != nil {
		return nil, apperrors.ModifyErr(*apperrors.ErrInternal, fmt.Errorf("logger init failed: %w", err))
	}

	s := &Server{
		serv:   service,
		logger: lg,
		httpServer: &http.Server{
			Addr:         ":" + conf.Port,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
	s.httpServer.Handler = lg.Middleware(s.SetupRoutes())

	return s, nil
}

func (s *Server) Run() error {
	s.logger.LogInfo(fmt.Sprintf("Server starting on %s\n", s.httpServer.Addr))

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.LogInfo("Server shutdown")
	fmt.Println("Server shutdown")
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /create_event", s.handleCreateEvent)
	mux.HandleFunc("POST /update_event", s.handleUpdateEvent)
	mux.HandleFunc("POST /delete_event", s.handleDeleteEvent)
	mux.HandleFunc("GET /events_for_day", s.handleGetDayEvents)
	mux.HandleFunc("GET /events_for_week", s.handleGetWeekEvents)
	mux.HandleFunc("GET /events_for_month", s.handleGetMonthEvents)
	return mux
}
