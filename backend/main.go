package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leeryan2000/flashat/internal/genproto/authpb"
	"github.com/leeryan2000/flashat/server"
	"github.com/leeryan2000/flashat/transport"
	"github.com/leeryan2000/flashat/transport/grpcserver"
	"google.golang.org/grpc"
)

func main() {
	// Structured JSON logging — output goes to stdout, Docker captures it
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	// ctx is cancelled the moment SIGINT/SIGTERM arrives — passed down to
	// MessageWorker so it can stop cleanly instead of being cut off mid-shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s, err := server.StartServer(ctx)
	if err != nil {
		log.Fatal("failed to start server: ", err)
	}

	router := transport.BuildRouter(s)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start HTTP server in background
	go func() {
		slog.Info("server listening", "port", 8080)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen error: ", err)
		}
	}()

	// Internal-only gRPC server — never proxied by nginx or published in
	// docker-compose, reachable only as backend:50051 on flashat-net.
	// Lets other services (currently just posts) reuse session/friendship
	// data without touching Redis/Postgres directly.
	grpcServer := grpc.NewServer()
	authpb.RegisterAuthInternalServer(grpcServer, &grpcserver.AuthServer{
		RedisClient:    s.RedisClient,
		FriendshipRepo: s.FriendshipRepo,
	})
	grpcLis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("grpc listen error: ", err)
	}
	go func() {
		slog.Info("grpc server listening", "port", 50051)
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatal("grpc serve error: ", err)
		}
	}()

	// Block until SIGINT or SIGTERM
	<-ctx.Done()
	slog.Info("shutting down...")

	// Give in-flight HTTP requests 10s to finish
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	httpServer.Shutdown(shutdownCtx)
	grpcServer.GracefulStop()

	// Let in-flight message deliveries finish (or hit their own timeout)
	// before pulling the connections they depend on out from under them.
	s.MessageWorker.WaitForShutdown()

	// Close connections cleanly
	s.RabbitMQClient.Close()
	s.RedisClient.Close()
	s.Pool.Close()

	slog.Info("shutdown complete")
}
