package db_test

import (
	"context"
	"log/slog"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/victor-nach/todo/internal-service/internal/db"
	"github.com/victor-nach/todo/internal-service/pkg/utils"
)

var dbURL string

func TestMain(m *testing.M) {
	var cleanup func()
	dbURL, cleanup = utils.CreateTestDB(&testing.T{})
	defer cleanup()

	m.Run()
}

func Test_MigrateUp(t *testing.T) {
	ctx := context.Background()
	storage, err := db.New(ctx, dbURL, slog.Default())
	require.NoError(t, err)
	defer storage.Close()

	err = storage.Migrate(ctx)
	require.NoError(t, err)
}