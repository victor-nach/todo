package handlers

import (
	"fmt"
	"strings"

	"github.com/microcosm-cc/bluemonday"

	"github.com/google/uuid"
	"github.com/victor-nach/todo/internal-service/internal/domain"
	pb "github.com/victor-nach/todo/proto/gen/go/todo"
)

func validateCreateTodoRequest(req *pb.CreateTodoRequest) error {
	var errs []string

	if req.GetTitle() == "" {
		errs = append(errs, "title is required")
	}
	title, err := validateField(req.GetTitle(), "title", 2, 50)
	if err != nil {
		errs = append(errs, err.Error())
	}
	req.Title = title

	// Validate if provided
	if req.GetDescription() != "" {
		title, err := validateField(req.GetTitle(), "title", 2, 50)
		if err != nil {
			errs = append(errs, err.Error())
		}
		req.Title = title
	}

	return joinErrors(errs)
}

func validateUpdateTodoRequest(req *pb.UpdateTodoRequest) (domain.UpdateParams, error) {
	var errs []string
	var params domain.UpdateParams

	if err := validateID(req.GetId()); err != nil {
		errs = append(errs, err.Error())
	}

	if req.GetTitle() == "" && req.GetDescription() == "" {
		return params, fmt.Errorf("at least one field must be provided")
	}

	title, err := validateField(req.GetTitle(), "title", 2, 100)
	if err != nil {
		errs = append(errs, err.Error())
	}
	params.Title = &title

	description, err := validateField(req.GetTitle(), "description", 2, 100)
	if err != nil {
		errs = append(errs, err.Error())
	}
	params.Description = &description

	return params, joinErrors(errs)
}

func validateID(id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}

	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("invalid id type provided")
	}

	return nil
}

func validateField(field, fieldName string, min, max int) (string, error) {
	field = sanitizeInput(field)
	if len(field) < min {
		return "", fmt.Errorf("%s length must be at least %d characters", fieldName, min)
	}
	if len(field) > max {
		return "", fmt.Errorf("%s length must not exceed %d characters", fieldName, max)
	}

	return field, nil
}

func sanitizeInput(input string) string {
	p := bluemonday.StrictPolicy()
	return p.Sanitize(input)
}

func joinErrors(errs []string) error {
	if len(errs) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errs, "; "))
	}

	return nil
}
