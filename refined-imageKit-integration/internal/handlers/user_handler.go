package handlers

import (
	"encoding/json"
	"imakit-practice/internal/models"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "image is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	image, err := h.Storage.UploadImage(file, header.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := models.User{
		Email: email,
		Image: *image,
		ID:    bson.NewObjectID(),
	}

	_, err = h.DB.Users.InsertOne(r.Context(), user)
	if err != nil {
		if deleteErr := h.Storage.DeleteImage(image.FileID); deleteErr != nil {
			log.Printf("rollback failed for fileID %s: %v", image.FileID, deleteErr)
		}
		log.Printf("db insert failed: %v", err)
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
