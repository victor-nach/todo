package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"log/slog"

	"github.com/uptrace/bunrouter"
	pb "github.com/victor-nach/todo/proto/gen/go/todo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -destination=./mocks/mock.go -package=mocks github.com/victor-nach/todo/api-gateway/internal/handlers GRPCClient
type GRPCClient interface {
	CreateTodo(ctx context.Context, in *pb.CreateTodoRequest, opts ...grpc.CallOption) (*pb.TodoResponse, error)
	GetTodo(ctx context.Context, in *pb.GetTodoRequest, opts ...grpc.CallOption) (*pb.TodoResponse, error)
	ListTodos(ctx context.Context, in *pb.ListTodosRequest, opts ...grpc.CallOption) (*pb.ListTodosResponse, error)
	UpdateTodo(ctx context.Context, in *pb.UpdateTodoRequest, opts ...grpc.CallOption) (*pb.TodoResponse, error)
	DeleteTodo(ctx context.Context, in *pb.DeleteTodoRequest, opts ...grpc.CallOption) (*pb.DeleteTodoResponse, error)
}

type Handler struct {
	grpcClient GRPCClient
	logger     *slog.Logger
}

func NewHandler(grpcClient GRPCClient, logger *slog.Logger) *Handler {
	return &Handler{
		grpcClient: grpcClient,
		logger:     logger,
	}
}


func (h *Handler) CreateTodo(w http.ResponseWriter, req bunrouter.Request) error {
	log := h.logger.With("method", "CreateTodo")
	ctx := req.Request.Context()

	var cr createReq
	if err := json.NewDecoder(req.Request.Body).Decode(&cr); err != nil {
		log.Error("failed to decode create todo request", "error", err)
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return err
	}

	grpcReq := &pb.CreateTodoRequest{
		Title:       cr.Title,
		Description: &cr.Description,
	}

	resp, err := h.grpcClient.CreateTodo(ctx, grpcReq)
	if err != nil {
		log.Error("failed to create todo", "error", err)
		respondWithErrorGRPC(w, err)
		return err
	}

	respondWithSuccess(w, mapTodo(resp.Todo))
	return nil
}

func (h *Handler) ListTodos(w http.ResponseWriter, req bunrouter.Request) error {
	log := h.logger.With("method", "ListTodos")
	ctx := req.Request.Context()

	resp, err := h.grpcClient.ListTodos(ctx, &pb.ListTodosRequest{})
	if err != nil {
		log.Error("failed to list todos", "error", err)
		respondWithErrorGRPC(w, err)
		return err
	}

	log.Info("Successfully listed todos", "count", len(resp.Todos))
	respondWithSuccess(w, mapTodos(resp.Todos))
	return nil
}

func (h *Handler) GetTodo(w http.ResponseWriter, req bunrouter.Request) error {
	log := h.logger.With("method", "GetTodo")
	ctx := req.Request.Context()

	id := req.Param("id")
	if id == "" {
		log.Error("missing id parameter")
		respondWithError(w, http.StatusBadRequest, "missing id parameter")
		return errors.New("missing id parameter")
	}

	resp, err := h.grpcClient.GetTodo(ctx, &pb.GetTodoRequest{Id: id})
	if err != nil {
		log.Error("failed to get todo", "error", err)
		respondWithErrorGRPC(w, err)
		return err
	}

	log.Info("Successfully retrieved todo", "id", resp.Todo.Id)
	respondWithSuccess(w, mapTodo(resp.Todo))
	return nil
}

func (h *Handler) UpdateTodo(w http.ResponseWriter, req bunrouter.Request) error {
	log := h.logger.With("method", "UpdateTodo")
	ctx := req.Request.Context()

	id := req.Param("id")
	if id == "" {
		log.Error("missing id parameter")
		respondWithError(w, http.StatusBadRequest, "missing id parameter")
		return errors.New("missing id parameter")
	}

	var ur updateReq
	if err := json.NewDecoder(req.Request.Body).Decode(&ur); err != nil {
		log.Error("failed to decode update todo request", "error", err)
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return err
	}

	grpcReq := &pb.UpdateTodoRequest{
		Id:          id,
		Title:       ur.Title,
		Description: &ur.Description,
	}
	resp, err := h.grpcClient.UpdateTodo(ctx, grpcReq)
	if err != nil {
		log.Error("failed to update todo", "error", err)
		respondWithErrorGRPC(w, err)
		return err
	}

	log.Info("Successfully updated todo", "id", resp.Todo.Id)
	respondWithSuccess(w, mapTodo(resp.Todo))
	return nil
}

func (h *Handler) DeleteTodo(w http.ResponseWriter, req bunrouter.Request) error {
	log := h.logger.With("method", "DeleteTodo")
	ctx := req.Request.Context()

	id := req.Param("id")
	if id == "" {
		log.Error("missing id parameter")
		respondWithError(w, http.StatusBadRequest, "missing id parameter")
		return errors.New("missing id parameter")
	}

	grpcReq := &pb.DeleteTodoRequest{Id: id}
	resp, err := h.grpcClient.DeleteTodo(ctx, grpcReq)
	if err != nil {
		log.Error("failed to delete todo", "error", err)
		respondWithErrorGRPC(w, err)
		return err
	}

	respondWithSuccess(w, resp)
	return nil
}

func respondWithSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response := APIResponse{
		Status: "success",
		Data:   data,
	}
	json.NewEncoder(w).Encode(response)
}

func respondWithError(w http.ResponseWriter, httpStatus int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	response := APIResponse{
		Status: "error",
		Error:  errorMsg,
	}
	json.NewEncoder(w).Encode(response)
}

func respondWithErrorGRPC(w http.ResponseWriter, err error) {
	var httpStatus int
	var errMsg string
	if st, ok := status.FromError(err); ok {
		errMsg = st.Message()
		switch st.Code() {
		case codes.InvalidArgument:
			httpStatus = http.StatusBadRequest
		case codes.NotFound:
			httpStatus = http.StatusNotFound
		default:
			httpStatus = http.StatusInternalServerError
		}
	} else {
		httpStatus = http.StatusInternalServerError
		errMsg = "internal server error"
	}
	respondWithError(w, httpStatus, errMsg)
}

func mapTodo(pbTodo *pb.Todo) Todo {
	var createdAt, updatedAt time.Time
	if pbTodo.CreatedAt != nil {
		createdAt = pbTodo.CreatedAt.AsTime()
	}
	if pbTodo.UpdatedAt != nil {
		updatedAt = pbTodo.UpdatedAt.AsTime()
	}

	var description string
	if pbTodo.Description != nil {
		description = *pbTodo.Description
	}

	if updatedAt.IsZero() {
		return Todo{
			ID:          pbTodo.Id,
			Title:       pbTodo.Title,
			Description: description,
			CreatedAt:   createdAt,
		}
	}
	return Todo{
		ID:          pbTodo.Id,
		Title:       pbTodo.Title,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   &updatedAt,
	}
}

func mapTodos(pbTodos []*pb.Todo) []Todo {
	todos := make([]Todo, 0, len(pbTodos))
	for _, pbTodo := range pbTodos {
		todos = append(todos, mapTodo(pbTodo))
	}
	return todos
}