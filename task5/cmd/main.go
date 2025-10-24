package main

import (
	"context"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pozedorum/WB_project_4/task5/internal/config"
	"github.com/pozedorum/WB_project_4/task5/internal/repository/postgres"
	"github.com/pozedorum/WB_project_4/task5/internal/server"
	"github.com/pozedorum/WB_project_4/task5/internal/service"
	"github.com/pozedorum/wbf/dbpg"
	"github.com/pozedorum/wbf/ginext"
	"github.com/pozedorum/wbf/zlog"
	"github.com/rs/zerolog"
)

func main() {
	zlog.Init()
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	// Настройки памяти для профилирования
	runtime.MemProfileRate = 1
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)

	startPProfServer()
	// Отключаем логирование Gin, для нагрузочного тестирования
	gin.SetMode(gin.ReleaseMode)
	cfg := config.Load()
	zlog.Logger.Info().Interface("config", cfg).Msg("Configuration loaded")

	opts := &dbpg.Options{
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	}

	pgRepo, err := postgres.NewShortURLRepositoryWithDB(cfg.Database.GetDSN(), []string{}, opts)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	defer func() {
		if pgRepo != nil {
			pgRepo.Close()
		}
		zlog.Logger.Info().Msg("PostgreSQL connection closed")
	}()

	shortURLService := service.New(pgRepo)
	server := server.New(shortURLService)
	router := ginext.New()
	router.LoadHTMLGlob("internal/frontend/templates/*.html")
	apiGroup := router.Group("")
	pprofGroup := router.Group("/debug/pprof")
	server.SetupRoutes(apiGroup)
	server.EnableProfiling(pprofGroup)

	// Запуск HTTP сервера в горутине
	serverAddr := ":" + cfg.Server.Port
	httpServer := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	go func() {
		zlog.Logger.Info().Str("address", serverAddr).Msg("Starting HTTP server")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zlog.Logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zlog.Logger.Info().Msg("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Останавливаем HTTP сервер
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		zlog.Logger.Error().Err(err).Msg("HTTP server shutdown error")
	} else {
		zlog.Logger.Info().Msg("HTTP server stopped gracefully")
	}
}

func startPProfServer() {
	// Создаем отдельный mux и РУЧНО регистрируем pprof handlers
	pprofMux := http.NewServeMux()

	// Регистрируем ВСЕ pprof handlers вручную
	pprofMux.HandleFunc("/debug/pprof/", pprof.Index)
	pprofMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	pprofMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	pprofMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	pprofMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)

	// Регистрируем отдельные профили
	pprofMux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	pprofMux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	pprofMux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	pprofMux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	pprofMux.Handle("/debug/pprof/block", pprof.Handler("block"))
	pprofMux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))

	server := &http.Server{
		Addr:    ":6060",
		Handler: pprofMux,
	}

	go func() {
		log.Printf("Starting pprof server on :6060")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("pprof server error: %v", err)
		}
	}()
}
