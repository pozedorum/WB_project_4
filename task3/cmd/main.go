package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	di "github.com/pozedorum/WB_project_4/task3/internal/DI"
	"github.com/pozedorum/WB_project_4/task3/pkg/config"
)

func main() {
	cfg := config.Load()
	aplicationContainer, err := di.NewContainer(cfg)
	if err != nil {
		fmt.Println("error while loading container: ", err)
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := aplicationContainer.Start(); err != nil {
			fmt.Println(err)
			// zlog.Logger.Fatal().Err(err).Msg("Failed to start server")
			return
		}
	}()

	<-quit
	fmt.Println("Shutting down server...")

	if err := aplicationContainer.Shutdown(); err != nil {
		// zlog.Logger.Fatal().Err(err).Msg("Forced shutdown")
		fmt.Println(err)
	}
}
