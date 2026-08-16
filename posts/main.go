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

	"github.com/leeryan2000/flashat-posts/server"
	"github.com/leeryan2000/flashat-posts/transport"
)

func main() {
	// Structured JSON logging — output goes to stdout, Docker captures it
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s, err := server.StartServer(ctx)
	if err != nil {
		log.Fatal("failed to start server: ", err)
	}

	router := transport.BuildRouter(s)
	httpServer := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	go func() {
		slog.Info("posts server listening", "port", 8081)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen error: ", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	httpServer.Shutdown(shutdownCtx)

	s.AuthConn.Close()
	s.Pool.Close()

	slog.Info("shutdown complete")
}
