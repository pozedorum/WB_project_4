// Package server описывает серверную часть сервиса
package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pozedorum/WB_project_4/task4/pkg/analyzer"
)

type Server struct {
	router   *gin.Engine
	analyzer *analyzer.Analyzer
}

const (
	pathHTML   = "internal/frontend/templates/*"
	pathStatic = "internal/frontend/static"
)

func NewServer() *Server {
	server := &Server{
		router:   gin.Default(),
		analyzer: analyzer.NewAnalyzer(),
	}

	server.SetupRoutes()
	return server
}

func (serv *Server) Run(addr string) error {
	return serv.router.Run(addr)
}

func (serv *Server) SetupRoutes() {
	serv.router.LoadHTMLGlob(pathHTML)
	serv.router.Static("/static", pathStatic)

	// Роуты API
	api := serv.router.Group("/api")
	{
		api.GET("/allocations", serv.AllocationsHandler)
		api.GET("/gc", serv.GCHandler)
		api.GET("/memory", serv.MemoryHandler)
		api.GET("/system", serv.SystemHandler)
		api.GET("/alloc", serv.MemoryAllocHandler)
	}

	// Роуты страниц
	serv.router.GET("/", serv.FrontendHandler)
	serv.router.GET("/metrics", serv.PrometheusMetricsHandler)
}

// FrontendHandler отображает главную страницу с кнопками
func (serv *Server) FrontendHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Go Metrics Dashboard",
	})
}
