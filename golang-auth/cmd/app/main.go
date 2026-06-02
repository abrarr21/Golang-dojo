package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abrarr21/golang-auth/internal/config"
	"github.com/abrarr21/golang-auth/internal/database"
	"github.com/abrarr21/golang-auth/internal/routes"
)

func main() {
	cfg := config.Load()
	db := database.ConnectDB(&cfg.Database)
	defer db.Disconnect()

	router := routes.RegisterAllRoutes(db, cfg)

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Server listening on port: ", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("server failed: %v", err)
		}
	}()

	sig := <-quit
	log.Println("Server shutdown signal received: ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("failed to shutdown gracefully, forcing server shutdown: %v", err)
	}

	log.Println("Server closed gracefully")
}
