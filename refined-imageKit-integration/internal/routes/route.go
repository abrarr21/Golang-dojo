package routes

import (
	"imakit-practice/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterAllRoutes(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	UserRoutes(r, h)
	return r
}
