package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leeryan2000/flashat/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getDSN() string {
	c := config.GetConfig() // Load config right when we need it
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DB_HOST, c.DB_USER, c.DB_PASS, c.DB_NAME, c.DB_PORT)
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
