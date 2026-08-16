package server

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leeryan2000/flashat-posts/config"
	"github.com/leeryan2000/flashat-posts/db"
	"github.com/leeryan2000/flashat-posts/grpcclient"
	"github.com/leeryan2000/flashat-posts/internal/genproto/authpb"
	repo "github.com/leeryan2000/flashat-posts/repo/impl"
	"google.golang.org/grpc"
)

type Server struct {
	Pool       *pgxpool.Pool
	AuthConn   *grpc.ClientConn
	AuthClient authpb.AuthInternalClient

	PostRepo *repo.PgxPostRepo
}

func StartServer(ctx context.Context) (*Server, error) {
	s := &Server{}

	cfg := config.LoadConfig()

	if err := db.EnsureDatabase(ctx, cfg); err != nil {
		return nil, err
	}

	if err := db.RunMigrations(cfg); err != nil {
		return nil, err
	}

	pool, err := db.NewPgxPool(cfg)
	if err != nil {
		return nil, err
	}
	s.Pool = pool

	authClient, authConn, err := grpcclient.NewAuthClient(cfg.AUTH_GRPC_ADDR)
	if err != nil {
		return nil, err
	}
	s.AuthClient = authClient
	s.AuthConn = authConn

	s.PostRepo = &repo.PgxPostRepo{Pool: s.Pool}

	return s, nil
}
