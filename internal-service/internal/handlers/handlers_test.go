package handlers_test

import (
	"context"
	"testing"
	"time"

	"log/slog"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/victor-nach/todo/internal-service/internal/domain"
	"github.com/victor-nach/todo/internal-service/internal/handlers"
	"github.com/victor-nach/todo/internal-service/internal/handlers/mocks"
	pb "github.com/victor-nach/todo/proto/gen/go/todo"
)

func TestHandler_CreateTodo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepo(ctrl)
	h := handlers.New(slog.Default(), mockRepo)

	req := &pb.CreateTodoRequest{
		Title:       "Test Todo",
		Description: String("Test Description"),
	}

	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, todo *domain.Todo) error {
			require.Equal(t, req.Title, todo.Title)
			return nil
		},
	).Times(1)

	resp, err := h.CreateTodo(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, req.Title, resp.Todo.Title)
}

func TestHandler_GetTodo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepo(ctrl)
	h := handlers.New(slog.Default(), mockRepo)

	id := uuid.NewString()
	req := &pb.GetTodoRequest{Id: id}

	expectedTodo := domain.Todo{
		ID:          id,
		Title:       "Test Todo",
		Description: "Test Description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockRepo.EXPECT().Get(gomock.Any(), id).Return(expectedTodo, nil).Times(1)

	resp, err := h.GetTodo(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, expectedTodo.Title, resp.Todo.Title)
	require.Equal(t, expectedTodo.Description, *resp.Todo.Description)
}

func TestHandler_ListTodos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepo(ctrl)
	h := handlers.New(slog.Default(), mockRepo)

	todos := []domain.Todo{
		{
			ID:          uuid.NewString(),
			Title:       "Todo 1",
			Description: "Description 1",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.NewString(),
			Title:       "Todo 2",
			Description: "Description 2",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	mockRepo.EXPECT().List(gomock.Any()).Return(todos, nil).Times(1)

	req := &pb.ListTodosRequest{}
	resp, err := h.ListTodos(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Todos, 2)
	require.Equal(t, todos[0].Title, resp.Todos[0].Title)
	require.Equal(t, todos[1].Title, resp.Todos[1].Title)
}

func String(s string) *string { return &s }
