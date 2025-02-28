package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"log/slog"

	"github.com/victor-nach/todo/internal-service/internal/config"
	"github.com/victor-nach/todo/internal-service/internal/db"
	"github.com/victor-nach/todo/internal-service/internal/handlers"
	"github.com/victor-nach/todo/internal-service/internal/repo"
	"github.com/victor-nach/todo/internal-service/internal/server"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	logger := slog.Default().With("app_env", cfg.AppEnv)

	storage, err := db.New(ctx, cfg.DBUrl, logger)
	if err != nil {
		logger.Error("failed to initialize repository", "err", err)
		os.Exit(1)
	}
	err = storage.Migrate(ctx)
	if err != nil {
		logger.Error("failed to migrate db", "err", err)
		os.Exit(1)
	}
	todoRepo := repo.New(storage.DB())

	grpcHandler := handlers.New(logger, todoRepo)

	grpcServer := server.New(grpcHandler, logger)

	// Start the server in a goroutine
	go func() {
		logger.Info("starting gRPC server", "port", cfg.Port)
		if err := grpcServer.Serve(cfg.Port); err != nil {
			logger.Error("gRPC server error", "err", err)
		}
	}()

	// Graceful shutdown: wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.Info("shutting down gRPC server...")

	// Allow some time for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	grpcServer.GracefulStop(shutdownCtx)

	logger.Info("gRPC server gracefully stopped")
}
