package shared

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/lib/pq"
)

// PostgresContainer wraps testcontainer PostgreSQL instance
type PostgresContainer struct {
	Container testcontainers.Container
	DB        *sql.DB
	ConnStr   string
}

// StartPostgresContainer starts a PostgreSQL testcontainer
func StartPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Wait for database to be ready
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresContainer{
		Container: pgContainer,
		DB:        db,
		ConnStr:   connStr,
	}, nil
}

// Terminate stops and removes the container
func (c *PostgresContainer) Terminate(ctx context.Context) error {
	if c.DB != nil {
		c.DB.Close()
	}
	if c.Container != nil {
		return c.Container.Terminate(ctx)
	}
	return nil
}

// LoadSchema loads SQL schema from file
func LoadSchema(ctx context.Context, db *sql.DB, schemaFile string) error {
	// Check if file exists at relative path first (from test directory)
	schema, err := os.ReadFile(schemaFile)
	if err != nil {
		// Try from project root
		schemaPath := filepath.Join("tests", "integration", "sqlboiler", schemaFile)
		schema, err = os.ReadFile(schemaPath)
		if err != nil {
			return fmt.Errorf("failed to read schema file: %w", err)
		}
	}

	if _, err := db.ExecContext(ctx, string(schema)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}
