package db

import (
	"context"
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leeryan2000/flashat/config"
)

func getDSN(c config.Configuration) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		c.DB_HOST, c.DB_USER, c.DB_PASS, c.DB_NAME, c.DB_PORT,
	)
}

func getURL(c config.Configuration) string {
	if c.GO_ENV == "development" {
		return fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			c.DB_USER, c.DB_PASS, c.DB_HOST, c.DB_PORT, c.DB_NAME,
		)
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		c.DB_USER, c.DB_PASS, c.DB_HOST, c.DB_PORT, c.DB_NAME,
	)
}

//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunMigrations(c config.Configuration) error {
	d, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("failed to load embedded migrations: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, getURL(c))
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrate up: %w", err)
	}

	log.Println("Database migrations applied successfully!")
	return nil
}

func NewPgxPool(c config.Configuration) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(getDSN(c))
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
