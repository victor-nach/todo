package handlers_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/uptrace/bunrouter"
	"github.com/victor-nach/todo/api-gateway/internal/handlers"
	"github.com/victor-nach/todo/api-gateway/internal/handlers/mocks"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/victor-nach/todo/proto/gen/go/todo"
)

func makeBunrouterRequest(r *http.Request) bunrouter.Request {
	return bunrouter.Request{
		Request: r,
	}
}

func TestHandler_CreateTodo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGRPCClient := mocks.NewMockGRPCClient(ctrl)

	h := handlers.NewHandler(mockGRPCClient, slog.Default())

	reqBody := `{"title": "Test Todo", "description": "Test Description"}`
	req := httptest.NewRequest("POST", "/todos", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	br := makeBunrouterRequest(req)

	now := time.Now()
	pbTodo := &pb.Todo{
		Id:          "12345",
		Title:       "Test Todo",
		Description: String("Test Description"),
		CreatedAt:   timestamppb.New(now),
		UpdatedAt:   timestamppb.New(now),
	}
	grpcResp := &pb.TodoResponse{
		Todo: pbTodo,
	}

	mockGRPCClient.
		EXPECT().
		CreateTodo(gomock.Any(), gomock.AssignableToTypeOf(&pb.CreateTodoRequest{}), gomock.Any()).
		DoAndReturn(func(ctx context.Context, req *pb.CreateTodoRequest, opts ...grpc.CallOption) (*pb.TodoResponse, error) {
			require.Equal(t, "Test Todo", req.Title)
			require.Equal(t, "Test Description", *req.Description)
			return grpcResp, nil
		}).
		Times(1)

	err := h.CreateTodo(rr, br)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, rr.Code)

	var apiResp handlers.APIResponse
	err = json.Unmarshal(rr.Body.Bytes(), &apiResp)
	require.NoError(t, err)
	require.Equal(t, "success", apiResp.Status)

	dataMap, ok := apiResp.Data.(map[string]interface{})
	require.True(t, ok, "expected Data to be a map")

	require.Equal(t, "12345", dataMap["id"])
	require.Equal(t, "Test Todo", dataMap["title"])
	require.Equal(t, "Test Description", dataMap["description"])
	require.NotEmpty(t, dataMap["created_at"])
}

func TestHandler_ListTodos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockGRPCClient(ctrl)
	h := handlers.NewHandler(mockClient, slog.Default())

	req := httptest.NewRequest("GET", "/todos", nil)
	rr := httptest.NewRecorder()
	br := makeBunrouterRequest(req)

	now := time.Now()
	pbTodo1 := &pb.Todo{
		Id:          "1",
		Title:       "Todo 1",
		Description: String("Description 1"),
		CreatedAt:   timestamppb.New(now),
		UpdatedAt:   timestamppb.New(now),
	}
	pbTodo2 := &pb.Todo{
		Id:          "2",
		Title:       "Todo 2",
		Description: String("Description 2"),
		CreatedAt:   timestamppb.New(now),
	}
	grpcResp := &pb.ListTodosResponse{
		Todos: []*pb.Todo{pbTodo1, pbTodo2},
	}

	mockClient.
		EXPECT().
		ListTodos(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(grpcResp, nil).
		Times(1)

	err := h.ListTodos(rr, br)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rr.Code)

	var apiResp handlers.APIResponse
	err = json.Unmarshal(rr.Body.Bytes(), &apiResp)
	require.NoError(t, err)
	require.Equal(t, "success", apiResp.Status)

	dataSlice, ok := apiResp.Data.([]interface{})
	require.True(t, ok)
	require.Len(t, dataSlice, 2)

	todo1, ok := dataSlice[0].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, "1", todo1["id"])
	require.Equal(t, "Todo 1", todo1["title"])
	require.Equal(t, "Description 1", todo1["description"])
	require.NotEmpty(t, todo1["created_at"])

	todo2, ok := dataSlice[1].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, "2", todo2["id"])
	require.Equal(t, "Todo 2", todo2["title"])
	require.Equal(t, "Description 2", todo2["description"])
	require.NotEmpty(t, todo2["created_at"])
}

func TestHandler_DeleteTodo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockGRPCClient(ctrl)
	h := handlers.NewHandler(mockClient, slog.Default())

	router := bunrouter.New()

	router.DELETE("/todos/:id", h.DeleteTodo)

	req := httptest.NewRequest(http.MethodDelete, "/todos/12345", nil)

	rr := httptest.NewRecorder()

	grpcResp := &pb.DeleteTodoResponse{Success: true}

	mockClient.
		EXPECT().
		DeleteTodo(gomock.Any(), &pb.DeleteTodoRequest{Id: "12345"}, gomock.Any()).
		Return(grpcResp, nil).
		Times(1)

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	var apiResp handlers.APIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &apiResp)
	require.NoError(t, err)
	require.Equal(t, "success", apiResp.Status)

	dataMap, ok := apiResp.Data.(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, true, dataMap["success"])
}

func TestHandler_GetTodo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockGRPCClient(ctrl)
	h := handlers.NewHandler(mockClient, slog.Default())

	router := bunrouter.New()

	router.GET("/todos/:id", h.GetTodo)

	req := httptest.NewRequest(http.MethodGet, "/todos/12345", nil)

	rr := httptest.NewRecorder()

	now := time.Now()
	pbTodo := &pb.Todo{
		Id:          "12345",
		Title:       "Test Todo",
		Description: String("Test Description"),
		CreatedAt:   timestamppb.New(now),
		UpdatedAt:   timestamppb.New(now),
	}
	grpcResp := &pb.TodoResponse{Todo: pbTodo}

	mockClient.
		EXPECT().GetTodo(gomock.Any(), &pb.GetTodoRequest{Id: "12345"}, gomock.Any()).
		Return(grpcResp, nil).Times(1)

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var apiResp handlers.APIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &apiResp)
	require.NoError(t, err)
	require.Equal(t, "success", apiResp.Status)

}

func String(s string) *string { return &s }
