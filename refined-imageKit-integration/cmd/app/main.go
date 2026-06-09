package main

import (
	"context"
	"imakit-practice/internal/config"
	"imakit-practice/internal/database"
	"imakit-practice/internal/handlers"
	"imakit-practice/internal/routes"
	"imakit-practice/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()
	db := database.ConnectDB(&cfg.Database)
	defer db.Disconnect()

	imageStorage := storage.NewImageKitStorage(&cfg.ImageKit)

	h := handlers.New(db, imageStorage)

	router := routes.RegisterAllRoutes(h)

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("server listening on port: ", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println("failed to start server")
		}
	}()

	sig := <-quit
	log.Println("shutdown signal recieved: ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("failed to close server cleanly, forcing server shutdown: %v", err)
	}

	log.Println("Server closed gracefully")

}
