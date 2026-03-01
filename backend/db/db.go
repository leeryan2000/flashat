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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var c = config.GetConfig() // Load config right when we need it

func getDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		c.DB_HOST, c.DB_USER, c.DB_PASS, c.DB_NAME, c.DB_PORT,
	)
}

func getURL() string {
	if c.GO_ENV == "development" {
		return fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			c.DB_USER, c.DB_PASS, c.DB_HOST, c.DB_PORT, c.DB_NAME,
		)
	} else {
		return fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			c.DB_USER, c.DB_PASS, c.DB_HOST, c.DB_PORT, c.DB_NAME,
		)
	}
}

//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunMigrations() error {

	// 1. Create the source driver from the embedded files
	d, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("failed to load embedded migrations: %v", err)
	}
	// 2. Create the migrate instance
	// Note: We use the "iofs" driver name here
	m, err := migrate.NewWithSourceInstance(
		"iofs", d, getURL())

	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}

	// 3. Run the "Up" migration
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrate up: %v", err)
	}

	log.Println("Database migrations applied successfully!")
	return nil
}

func InitDB() (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewPgxPool() (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// set the config of pgx here
	config, err := pgxpool.ParseConfig(getDSN())
	if err != nil {
		return nil, err
	}

	// Connect
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
