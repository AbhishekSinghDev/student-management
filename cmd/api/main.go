package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AbhishekSinghDev/student-management/internal/config"
	"github.com/AbhishekSinghDev/student-management/internal/http/handlers/health"
	"github.com/AbhishekSinghDev/student-management/internal/http/handlers/student"
	"github.com/AbhishekSinghDev/student-management/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	storage, err := sqlite.New(cfg)
	if err != nil {
		logger.Error("failed to initialize storage", "err", err)
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health.New())
	mux.HandleFunc("POST /api/student", student.New(storage))

	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("server starting", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "err", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown — give in-flight requests 30s to finish
	logger.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("forced shutdown", "err", err)
	}

	logger.Info("server stopped")
}
