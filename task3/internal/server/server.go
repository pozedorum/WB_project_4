package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pozedorum/WB_project_4/task3/internal/interfaces"
)

type EventServer struct {
	server *http.Server
	router *gin.Engine
	serv   interfaces.Service
	logger interfaces.Logger
}

func NewEventServer(port string, service interfaces.Service, logger interfaces.Logger) *EventServer {
	router := gin.Default()

	s := &EventServer{
		serv:   service,
		logger: logger,
		server: &http.Server{
			Addr:    ":" + port,
			Handler: router,
		},
		router: router,
	}

	s.setupRoutes()
	return s
}

func (s *EventServer) setupRoutes() {
	// CRUD операции для событий
	s.router.POST("/create_event", s.handleCreateEvent)
	s.router.POST("/update_event", s.handleUpdateEvent)
	s.router.POST("/delete_event", s.handleDeleteEvent)
	s.router.GET("/events_for_day", s.handleGetDayEvents)
	s.router.GET("/events_for_week", s.handleGetWeekEvents)
	s.router.GET("/events_for_month", s.handleGetMonthEvents)
}

func (s *EventServer) Start() error {
	s.logger.Info("SERVER_START", "Gin server starting", "addr", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *EventServer) Shutdown(ctx context.Context) error {
	s.logger.Info("SERVER_SHUTDOWN", "Initiating server shutdown")
	return s.server.Shutdown(ctx)
}
