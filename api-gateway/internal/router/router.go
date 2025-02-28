package router

import (
	"net/http"

	"github.com/uptrace/bunrouter"
	"github.com/victor-nach/todo/api-gateway/internal/handlers"
)

func NewRouter(h *handlers.Handler) http.Handler {
	r := bunrouter.New()

	r.GET("/todos", h.ListTodos)
	r.POST("/todos", h.CreateTodo)
	r.GET("/todos/:id", h.GetTodo)
	r.PATCH("/todos/:id", h.UpdateTodo)
	r.DELETE("/todos/:id", h.DeleteTodo)

	return r
}