package domain

import (
	"time"
)

type Todo struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateParams struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

