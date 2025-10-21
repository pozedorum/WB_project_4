package server

import (
	"github.com/pozedorum/WB_project_4/task4/internal/analyzer"
)

type Server struct {
	analyzer *analyzer.Analyzer
}

func NewServer() *Server {
	return &Server{
		analyzer: analyzer.NewAnalyzer(),
	}
}
