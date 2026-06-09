package routes

import (
	"imakit-practice/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func UserRoutes(r chi.Router, h *handlers.Handler) {

	r.Post("/img", h.CreateUser)
}
