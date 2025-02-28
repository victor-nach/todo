package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"log/slog"

	"github.com/victor-nach/todo/api-gateway/internal/config"
	"github.com/victor-nach/todo/api-gateway/internal/handlers"
	"github.com/victor-nach/todo/api-gateway/internal/router"

	pb "github.com/victor-nach/todo/proto/gen/go/todo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	logger := slog.Default()
	logger = logger.With("app_env", cfg.AppEnv)

	grpcConn, err := grpc.NewClient(
		cfg.GRPCServerURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("failed to create grpc client", "err", err)
		os.Exit(1)
	}
	grpcClient := pb.NewToDoServiceClient(grpcConn)

	h := handlers.NewHandler(grpcClient, logger)

	r := router.NewRouter(h)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Start server
	go func() {
		logger.Info("starting http server", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "err", err)
		}
	}()

	// Graceful shutdown: wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.Info("shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Error("server shutdown error", "err", err)
	}

	logger.Info("server gracefully stopped")
}
