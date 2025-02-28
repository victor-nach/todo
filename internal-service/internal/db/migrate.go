package db

import (
	"context"
	"embed"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun/migrate"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func (s *Storage) Migrate(ctx context.Context) error {
	migrations := migrate.NewMigrations()

	if err := migrations.Discover(migrationsFS); err != nil {
		return fmt.Errorf("failed to discover migrations: %w", err)
	}

	migrator := migrate.NewMigrator(s.db, migrations)

	// Initialize the bun_migrations table if it doesn't exist
	if err := migrator.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize migrations table: %w", err)
	}

	// Run all pending migrations
	group, err := migrator.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if group.IsZero() {
		log.Println("no new migrations to run")
		return nil
	}

	log.Printf("successfully ran %d migrations\n", len(group.Migrations))
	return nil
}
