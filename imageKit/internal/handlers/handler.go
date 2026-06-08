package handlers

import (
	"net/http"

	"github.com/abrarr21/auth-practice-3/internal/config"
	"github.com/abrarr21/auth-practice-3/internal/database"
)

type Handler struct {
	DB  *database.Database
	Cfg *config.Config
}

func New(db *database.Database, cfg *config.Config) *Handler {
	return &Handler{
		DB:  db,
		Cfg: cfg,
	}
}

func (h *Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Server running perfectly ✅"))
}
