package db

import (
	"context"
	"fmt"
	"time"

	"github.com/leeryan2000/flashat/config"
	"github.com/redis/go-redis/v9"
)

const sessionPrefix = "session:"
const sessionTTL = 3 * 24 * time.Hour

type RedisClient struct {
	Conn *redis.Client
}

// NewClient creates and verifies a Redis client.
func NewClient(cfg config.Configuration) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.REDIS_ADDR,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 2,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis at %s: %w", cfg.REDIS_ADDR, err)
	}

	return &RedisClient{
		Conn: rdb,
	}, nil
}

// Close shuts down the Redis connection cleanly.
func (c *RedisClient) Close() error {
	if c == nil || c.Conn == nil {
		return nil
	}
	return c.Conn.Close()
}

func (c *RedisClient) SetSession(ctx context.Context, sessionID string, uid string) error {
	return c.Conn.Set(ctx, sessionPrefix+sessionID, uid, sessionTTL).Err()
}

func (c *RedisClient) GetSession(ctx context.Context, sessionID string) (string, error) {
	return c.Conn.Get(ctx, sessionPrefix+sessionID).Result()
}

func (c *RedisClient) DeleteSession(ctx context.Context, sessionID string) error {
	return c.Conn.Del(ctx, sessionPrefix+sessionID).Err()
}
