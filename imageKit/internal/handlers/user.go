package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/abrarr21/auth-practice-3/internal/models"
	"github.com/abrarr21/auth-practice-3/internal/services"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	var existing models.User

	err = h.DB.Users.
		FindOne(
			context.Background(),
			bson.M{"email": email},
		).
		Decode(&existing)

	if err == nil {
		http.Error(w, "user already exists", http.StatusConflict)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "image required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	image, err := services.UploadImage(
		file,
		header.Filename,
	)

	if err != nil {
		log.Printf("failed to upload image: %v", err)
		http.Error(w, "failed to upload image", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
		Image:    *image,
	}

	result, err := h.DB.
		Users.InsertOne(
		context.Background(),
		user,
	)

	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	user.ID = result.InsertedID.(bson.ObjectID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(user)
}
