package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leeryan2000/flashat/server"
	"github.com/leeryan2000/flashat/transport"
)

func main() {
	s, err := server.StartServer()
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
		log.Println("server listening on :8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen error: ", err)
		}
	}()

	// Block until SIGINT or SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	// Give in-flight HTTP requests 10s to finish
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	httpServer.Shutdown(ctx)

	// Close connections cleanly
	s.RabbitMQClient.Close()
	s.RedisClient.Close()
	s.Pool.Close()

	log.Println("shutdown complete")
}
