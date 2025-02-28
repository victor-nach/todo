package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/uptrace/bun"
	"github.com/victor-nach/todo/internal-service/internal/domain"
)

type repo struct {
	db *bun.DB
}

func New(db *bun.DB) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, todo *domain.Todo) error {
	_, err := r.db.NewInsert().Model(todo).Exec(ctx)
	return err
}

func (r *repo) Get(ctx context.Context, id string) (domain.Todo, error) {
	var todo domain.Todo
	err := r.db.NewSelect().Model(&todo).Where("id = ?", id).Scan(ctx)
	if err != nil {	
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Todo{}, domain.ErrTodoNotFound
		}

		return domain.Todo{}, err
	}

	return todo, err
}


func (r *repo) List(ctx context.Context) ([]domain.Todo, error) {
	var todos []domain.Todo
	err := r.db.NewSelect().Model(&todos).Order("created_at DESC").Scan(ctx)
	return todos, err
}


func (r *repo) Update(ctx context.Context, id string, updateParams domain.UpdateParams) (domain.Todo, error) {
	var todo domain.Todo

	query := r.db.NewUpdate().Model(&todo).Where("id = ?", id)
	if updateParams.Title != nil {
		query = query.Set("title = ?", *updateParams.Title)
	}
	if updateParams.Description != nil {
		query = query.Set("description = ?", *updateParams.Description)
	}
	query = query.Set("updated_at = ?", time.Now()).Returning("*")

	err := query.Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Todo{}, domain.ErrTodoNotFound
		}
		return domain.Todo{}, err
	}

	return todo, nil
}

func (r *repo) Delete(ctx context.Context, id string) error {
	res, err := r.db.NewDelete().Model(&domain.Todo{}).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrTodoNotFound
	}
	
	return nil
}


