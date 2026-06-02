package routes

import (
	"github.com/abrarr21/golang-auth/internal/handlers"
	"github.com/abrarr21/golang-auth/internal/middlewares"
	"github.com/go-chi/chi/v5"
)

func UserRoutes(r chi.Router, h *handlers.Handler) {
	r.Route("/auth/users", func(r chi.Router) {
		r.Post("/register", h.RegisterUser)
		r.Post("/login", h.Login)
	})

	r.Group(func(r chi.Router) {
		r.Use(middlewares.RequireAuth(h.Cfg.JWT.JWT_SECRET))

		r.Get("/getme", h.GetMe)
	})
}
