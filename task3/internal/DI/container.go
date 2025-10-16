package di

import "github.com/pozedorum/WB_project_4/task3/internal/interfaces"

type Container struct {
	server  interfaces.Server
	service interfaces.Service
}
