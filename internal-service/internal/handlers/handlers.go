package handlers

import (
	"context"
	"errors"
	"time"

	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/victor-nach/todo/internal-service/internal/domain"
	pb "github.com/victor-nach/todo/proto/gen/go/todo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	pb.UnimplementedToDoServiceServer
	logger *slog.Logger
	repo   repo
}

//go:generate mockgen -destination=./mocks/mock.go -package=mocks github.com/victor-nach/todo/internal-service/internal/handlers repo
type repo interface {
	Create(ctx context.Context, todo *domain.Todo) error
	Get(ctx context.Context, id string) (domain.Todo, error)
	List(ctx context.Context) ([]domain.Todo, error)
	Update(ctx context.Context, id string, updateParams domain.UpdateParams) (domain.Todo, error)
	Delete(ctx context.Context, id string) error
}

func New(logger *slog.Logger, todosRepo repo) *Handler {
	logger = logger.With("package", "handlers")

	return &Handler{
		logger: logger,
		repo:   todosRepo,
	}
}

func (h *Handler) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.TodoResponse, error) {
	log := h.logger.With("method", "CreateTodo").With("title", req.GetTitle())

	log.Info("Handling create todo request")

	if err := validateCreateTodoRequest(req); err != nil {
		log.Error("Invalid request", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "error validating request: %v", err)
	}

	todo := &domain.Todo{
		ID:          uuid.NewString(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		CreatedAt:   time.Now(),
	}

	if err := h.repo.Create(ctx, todo); err != nil {
		log.Error("Failed to create todo", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create todo: %v", err)
	}

	log.Info("Todo created successfully", "id", todo.ID)

	return &pb.TodoResponse{
		Todo: mapTodoModelToPb(*todo),
	}, nil
}

func (h *Handler) GetTodo(ctx context.Context, req *pb.GetTodoRequest) (*pb.TodoResponse, error) {
	log := h.logger.With("method", "GetTodo").With("id", req.GetId())

	log.Info("Handling get todo request")

	if err := validateID(req.GetId()); err != nil {
		log.Error("Invalid request", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "error validating request: %v", err)
	}

	todo, err := h.repo.Get(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, domain.ErrTodoNotFound) {
			return nil, status.Errorf(codes.NotFound, "failed to get Todo: %v", err)
		}

		log.Error("Failed to get todo", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get todo: %v", err)
	}

	log.Info("Todo retrieved successfully")

	return &pb.TodoResponse{
		Todo: mapTodoModelToPb(todo),
	}, nil
}

func (h *Handler) ListTodos(ctx context.Context, req *pb.ListTodosRequest) (*pb.ListTodosResponse, error) {
	log := h.logger.With("method", "ListTodos")

	log.Info("Handling list todos request")

	todos, err := h.repo.List(ctx)
	if err != nil {
		log.Error("Failed to list todos", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list todos: %v", err)
	}

	log.Info("Todos listed successfully")

	var pbTodos []*pb.Todo
	for _, todo := range todos {
		pbTodos = append(pbTodos, mapTodoModelToPb(todo))
	}
	
	return &pb.ListTodosResponse{
		Todos: pbTodos,
	}, nil

}

func (h *Handler) UpdateTodo(ctx context.Context, req *pb.UpdateTodoRequest) (*pb.TodoResponse, error) {
	log := h.logger.With("method", "UpdateTodo").With("id", req.GetId())

	log.Info("Handling update todo request")

	updateParams, err := validateUpdateTodoRequest(req)
	if  err != nil {
		log.Error("Invalid request", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "error validating request: %v", err)
	}

	updatedTodo, err := h.repo.Update(ctx, req.GetId(), updateParams)
	if err != nil {
		if errors.Is(err, domain.ErrTodoNotFound) {
			return nil, status.Errorf(codes.NotFound, "failed to update Todo: %v", err)
		}

		log.Error("Failed to update todo", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get todo: %v", err)
	}

	log.Info("Todo updated successfully")

	return &pb.TodoResponse{
		Todo: mapTodoModelToPb(updatedTodo),
	}, nil
}

func (h *Handler) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*pb.DeleteTodoResponse, error) {
	log := h.logger.With("method", "DeleteTodo").With("id", req.GetId())

	log.Info("Handling delete todo request")

	if err := validateID(req.GetId()); err != nil {
		log.Error("Invalid request", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "error validating request: %v", err)
	}

	if err := h.repo.Delete(ctx, req.GetId()); err != nil {
		if errors.Is(err, domain.ErrTodoNotFound) {
			return nil, status.Errorf(codes.NotFound, "failed to delete Todo: %v", err)
		}

		log.Error("Failed to delete todo", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to delete todo: %v", err)
	}

	log.Info("Todo deleted successfully")

	return &pb.DeleteTodoResponse{
		Success: true,
	}, nil
}

func mapTodoModelToPb(todo domain.Todo) *pb.Todo {
	return &pb.Todo{
		Id:          todo.ID,
		Title:       todo.Title,
		Description: &todo.Description,
		CreatedAt:   timestamppb.New(todo.CreatedAt),
		UpdatedAt:   timestamppb.New(todo.UpdatedAt),
	}
}
