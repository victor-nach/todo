package utils

import (
	"fmt"
	"log"
	"testing"
	"time"

	"database/sql"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

// CreateTestDB initializes an in-memory PostgreSQL database using dockertest
// returns the database url and and a cleanup function
func CreateTestDB(t *testing.T) (string, func()) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not create dockertest pool: %v", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120)

	var db *sql.DB
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Clean up function to stop the container and release resources
	cleanup := func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge the container: %v", err)
		}
	}

	return databaseUrl, cleanup
}
