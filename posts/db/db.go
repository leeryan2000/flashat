package db

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leeryan2000/flashat-posts/config"
)

func getDSN(c config.Configuration, dbname string) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		c.DB_HOST, c.DB_USER, c.DB_PASS, dbname, c.DB_PORT,
	)
}

func getURL(c config.Configuration, dbname string) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DB_USER, c.DB_PASS, c.DB_HOST, c.DB_PORT, dbname,
	)
}

// EnsureDatabase creates the posts service's own logical database
// (config.PostsDBName) within the monolith's Postgres instance, if it
// doesn't already exist. Needed because postgres:17-alpine's POSTGRES_DB
// env var only creates a database on first init of an empty volume — on
// an already-initialized EC2 volume it wouldn't retroactively create a
// second database, so this step does it explicitly and idempotently.
// Connects to the "postgres" maintenance database using the same admin
// credentials the monolith already uses (a dedicated least-privilege
// role is a reasonable future improvement, not required at this scale).
func EnsureDatabase(ctx context.Context, c config.Configuration) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, getDSN(c, "postgres"))
	if err != nil {
		return fmt.Errorf("failed to connect to maintenance database: %w", err)
	}
	defer conn.Close(ctx)

	var exists bool
	err = conn.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)`, config.PostsDBName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check for existing database: %w", err)
	}

	if exists {
		return nil
	}

	// Database names can't be parameterized in Postgres — safe here since
	// config.PostsDBName is a fixed constant, never user input.
	_, err = conn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE %s OWNER %s`, config.PostsDBName, c.DB_USER))
	if err != nil {
		return fmt.Errorf("failed to create database %s: %w", config.PostsDBName, err)
	}

	slog.Info("created posts database", "database", config.PostsDBName)
	return nil
}

//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunMigrations(c config.Configuration) error {
	d, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("failed to load embedded migrations: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, getURL(c, config.PostsDBName))
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrate up: %w", err)
	}

	slog.Info("posts database migrations applied successfully")
	return nil
}

func NewPgxPool(c config.Configuration) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(getDSN(c, config.PostsDBName))
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
