package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-financial-planning/backend/internal/app"
	"github.com/go-financial-planning/backend/internal/db"
	"github.com/go-financial-planning/backend/internal/handlers"
	"github.com/go-financial-planning/backend/internal/repository"
)

func main() {
	cfg := app.LoadConfig()

	database, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer database.Close()

	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	repo := repository.New(database)
	handler := handlers.New(cfg, repo)

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler.Router(),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("api listening on %s", cfg.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %v", err)
	}
}
