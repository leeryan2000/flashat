package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leeryan2000/flashat/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var configuration = config.GetConfig()

var dsn = "host=" + configuration.DB_DBHOST + " " +
	"user=" + configuration.DB_USER + " " +
	"password=" + configuration.DB_PASS + " " +
	"dbname=" + configuration.DB_DBNAME + " " +
	"port=" + configuration.DB_PORT + " " +
	"sslmode=disable"

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewPgxPool() (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)
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
