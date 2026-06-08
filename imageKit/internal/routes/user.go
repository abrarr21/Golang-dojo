package routes

import (
	"github.com/abrarr21/auth-practice-3/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func UserRoutes(r chi.Router, h *handlers.Handler) {
	r.Route("/auth/users", func(r chi.Router) {
		r.Post("/register", h.CreateUser)
	})
}
