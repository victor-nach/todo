package server

import (
	"context"
	"fmt"
	"net"

	pb "github.com/victor-nach/todo/proto/gen/go/todo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
)

type Server struct {
	grpcServer *grpc.Server
	logger     *slog.Logger
}

func New(handler pb.ToDoServiceServer, logger *slog.Logger) *Server {
	s := grpc.NewServer()
	pb.RegisterToDoServiceServer(s, handler)
	reflection.Register(s)
	return &Server{
		grpcServer: s,
		logger:     logger,
	}
}

// Serve starts the gRPC server on the provided port.
func (s *Server) Serve(port string) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		s.logger.Error("failed to listen", "err", err)
		return err
	}
	s.logger.Info("gRPC server is listening", "port", port)
	return s.grpcServer.Serve(listener)
}

// GracefulStop attempts a graceful stop within the provided timeout.
func (s *Server) GracefulStop(ctx context.Context) {
	stopped := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()
	select {
	case <-ctx.Done():
		s.logger.Error("graceful stop timeout, forcing stop")
		s.grpcServer.Stop()
	case <-stopped:
		s.logger.Info("graceful stop complete")
	}
}
