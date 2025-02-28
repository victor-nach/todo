package repo_test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/victor-nach/todo/internal-service/internal/db"
	"github.com/victor-nach/todo/internal-service/internal/domain"
	"github.com/victor-nach/todo/internal-service/internal/repo"
	"github.com/victor-nach/todo/internal-service/pkg/utils"
)

var storage *db.Storage

func TestMain(m *testing.M) {
	dbURL, cleanup := utils.CreateTestDB(&testing.T{})
	defer cleanup()

	ctx := context.Background()
	var err error
	storage, err = db.New(ctx, dbURL, slog.Default())
	if err != nil {
		fmt.Printf("failed to connect to db: %v", err)
		return
	}

	defer storage.Close()

	err = storage.Migrate(ctx)
	if err != nil {
		fmt.Printf("failed to migrate db: %v", err)
		return
	}

	m.Run()
}

func Test_Create(t *testing.T) {
	ctx := context.Background()
	repo := repo.New(storage.DB())

	todo := domain.Todo{
		ID:          uuid.NewString(),
		Title:       "test",
		Description: "test",
		CreatedAt:   time.Now(),
	}

	err := repo.Create(ctx, &todo)
	require.NoError(t, err)
}

func Test_Get(t *testing.T) {
	ctx := context.Background()
	repo := repo.New(storage.DB())

	t.Run("ok", func(t *testing.T) {
		todo := createTodo(1)[0]

		err := repo.Create(ctx, &todo)
		require.NoError(t, err)

		todo, err = repo.Get(ctx, todo.ID)
		require.NoError(t, err)
		require.NotEmpty(t, todo)
		assert.Equal(t, "test title 1", todo.Title)
		assert.Equal(t, "test description 1", todo.Description)
		assert.NotEmpty(t, todo.CreatedAt)
		assert.Equal(t, true, todo.UpdatedAt.IsZero())
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.Get(ctx, uuid.NewString())
		require.Error(t, err)
		assert.Equal(t, domain.ErrTodoNotFound, err)
	})
}

func Test_List(t *testing.T) {
	ctx := context.Background()
	repo := repo.New(storage.DB())

	todo := createTodo(3)
	for _, v := range todo {
		err := repo.Create(ctx, &v)
		require.NoError(t, err)
	}

	todos, err := repo.List(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, todos)
}

func Test_Update(t *testing.T) {
	ctx := context.Background()
	repo := repo.New(storage.DB())

	t.Run("ok", func(t *testing.T) {
		todo := createTodo(1)[0]

		err := repo.Create(ctx, &todo)
		require.NoError(t, err)

		updateParams := domain.UpdateParams{
			Title:       String("updated title"),
			Description: String("updated description"),
		}

		todo, err = repo.Update(ctx, todo.ID, updateParams)
		require.NoError(t, err)
		require.NotEmpty(t, todo)
		assert.Equal(t, "updated title", todo.Title)
		assert.Equal(t, "updated description", todo.Description)
		assert.NotEmpty(t, todo.CreatedAt)
		assert.NotEmpty(t, todo.UpdatedAt)
	})

	t.Run("not found", func(t *testing.T) {
		updateParams := domain.UpdateParams{
			Title: String("updated title")}

		_, err := repo.Update(ctx, uuid.NewString(), updateParams)
		require.Error(t, err)
		assert.Equal(t, domain.ErrTodoNotFound, err)
	})
}

func createTodo(n int) []domain.Todo {
	todos := make([]domain.Todo, n)
	for i := 0; i < n; i++ {
		todos[i] = domain.Todo{
			ID:          uuid.NewString(),
			Title:       fmt.Sprintf("test title %d", i+1),
			Description: fmt.Sprintf("test description %d", i+1),
			CreatedAt:   time.Now(),
		}
	}

	return todos
}

func String(s string) *string {
	return &s
}
