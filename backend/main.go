package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leeryan2000/flashat/server"
	"github.com/leeryan2000/flashat/transport"
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

	// Block until SIGINT or SIGTERM
	<-ctx.Done()
	slog.Info("shutting down...")

	// Give in-flight HTTP requests 10s to finish
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	httpServer.Shutdown(shutdownCtx)

	// Let in-flight message deliveries finish (or hit their own timeout)
	// before pulling the connections they depend on out from under them.
	s.MessageWorker.WaitForShutdown()

	// Close connections cleanly
	s.RabbitMQClient.Close()
	s.RedisClient.Close()
	s.Pool.Close()

	slog.Info("shutdown complete")
}
