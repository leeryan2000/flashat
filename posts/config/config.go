package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// PostsDBName is intentionally not an env var — the posts service always
// owns a separate logical database within the same Postgres instance the
// monolith uses, never the monolith's own DB_NAME.
const PostsDBName = "flashat_posts"

type Configuration struct {
	// Reuses the monolith's Postgres admin credentials/host — same
	// cluster, separate logical database (PostsDBName). See db/db.go's
	// EnsureDatabase for why a dedicated least-privilege role isn't set
	// up here yet.
	DB_USER string
	DB_PASS string
	DB_HOST string
	DB_PORT string

	// AUTH_GRPC_ADDR is the backend's internal gRPC listener, reachable
	// only on the flashat-net bridge network (e.g. "backend:50051").
	AUTH_GRPC_ADDR string
}

func LoadConfig() Configuration {
	envFile := ".env"

	if err := godotenv.Load(envFile); err != nil {
		slog.Warn("error loading env file, reading from system env only", "file", envFile)
	}

	return Configuration{
		DB_USER:        os.Getenv("DB_USER"),
		DB_PASS:        os.Getenv("DB_PASS"),
		DB_HOST:        os.Getenv("DB_HOST"),
		DB_PORT:        os.Getenv("DB_PORT"),
		AUTH_GRPC_ADDR: os.Getenv("AUTH_GRPC_ADDR"),
	}
}
