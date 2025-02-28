package db

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type Storage struct {
	db *bun.DB
}

// New creates new db connection with retry
func New(ctx context.Context, dsn string, logger *slog.Logger)  (*Storage, error) {
	log := logger.With("package", "db").With("dsn", dsn)

	log.Info("connecting to db")

	sqldb, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	for i := 0; i < 5; i++ {

		err = sqldb.PingContext(ctx)
		if err == nil {
			break
		}

		log.Error("error connecting to db", "error", err, "attempt", i+1)
		time.Sleep(5 * time.Second) // Retry every 5 seconds
	}

	if err != nil {
		log.Error("All db connect attempts failed", "error", err)
		return nil, err
	}
	db := bun.NewDB(sqldb, pgdialect.New())

	log.Info("Successfully connected to db")

	return &Storage{db: db}, nil
}

// DB returns the current db connection
func (s *Storage) DB() *bun.DB {
	return s.db
}

// Close closes the db connection
func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
