package routes

import (
	"github.com/abrarr21/auth-practice-3/internal/config"
	"github.com/abrarr21/auth-practice-3/internal/database"
	"github.com/abrarr21/auth-practice-3/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterAllRoutes(db *database.Database, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	h := handlers.New(db, cfg)

	r.Get("/", h.CheckHealth)

	UserRoutes(r, h)

	return r
}
