package handlers

import "time"

type (
	createReq struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	updateReq struct {
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
	}

	APIResponse struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Error  string      `json:"error,omitempty"`
	}

	Todo struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description,omitempty"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	}
)
