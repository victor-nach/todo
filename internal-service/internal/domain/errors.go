package domain

import "errors"

var ErrTodoNotFound = errors.New("todo not found, invalid id provided")