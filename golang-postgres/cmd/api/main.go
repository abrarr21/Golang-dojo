package main

import (
	"fmt"
	"golang-postgres/internal/config"
	"golang-postgres/internal/db"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()
	_, err := db.ConnectDB(cfg.DbURL)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running properly"))
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	fmt.Println("server started running on port: ", cfg.Port)
	fmt.Println("Database connected")
	if err := srv.ListenAndServe(); err != nil {
		srv.ErrorLog.Fatal("server failed")
	}
}
