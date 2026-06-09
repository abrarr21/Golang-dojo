package handlers

import (
	"imakit-practice/internal/database"
	"imakit-practice/internal/storage"
	"net/http"
)

type Handler struct {
	DB      *database.Database
	Storage storage.Storage
}

func New(db *database.Database, s storage.Storage) *Handler {
	return &Handler{
		DB:      db,
		Storage: s,
	}
}

func (h *Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Server running perfectly"))
}
