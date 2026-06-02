package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/abrarr21/golang-auth/internal/middlewares"
	"github.com/abrarr21/golang-auth/internal/models"
	"github.com/abrarr21/golang-auth/internal/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user *models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("failed to decode input request: %v", err)
		utils.ResponseJSON(w, http.StatusBadRequest, "invalid input", nil)
		return
	}

	if err := utils.Validator.Struct(user); err != nil {
		log.Printf("Validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash the Passwrod", http.StatusInternalServerError)
		return
	}

	user.ID = bson.NewObjectID()
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.ToLower(user.Email)
	user.Password = string(hashed)

	result, err := h.DB.Users.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			http.Error(w, "email already exists", http.StatusConflict)
			return
		}
		log.Printf("failed to insert user: %v", err)
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	log.Println("User Created: ", result.InsertedID)

	token, err := utils.GenerateToken(user.ID.Hex(), user.Email, h.Cfg.JWT.JWT_SECRET, h.Cfg.JWT.AccessTokenTTL)
	if err != nil {
		log.Printf("failed to generate token: %v", err)
		utils.ResponseJSON(w, http.StatusInternalServerError, "failed to generate token", nil)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    token,
		MaxAge:   15 * 60,
		Path:     "/",
		HttpOnly: true,
	})

	response := models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	if err := utils.ResponseJSON(w, http.StatusCreated, "user created successfully", response); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input *models.LoginInput

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("failed to decode request: %v", err)
		utils.ResponseJSON(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := utils.Validator.Struct(input); err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	input.Email = strings.ToLower(input.Email)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var user models.User

	err := h.DB.Users.FindOne(ctx, bson.D{{Key: "email", Value: input.Email}}).Decode(&user)
	if err != nil {
		utils.ResponseJSON(w, http.StatusUnauthorized, "invalid email or password", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		utils.ResponseJSON(w, http.StatusUnauthorized, "invalid email or password", nil)
		return
	}

	token, err := utils.GenerateToken(user.ID.Hex(), user.Email, h.Cfg.JWT.JWT_SECRET, h.Cfg.JWT.AccessTokenTTL)
	if err != nil {
		log.Printf("failed to generate token: %v", err)
		utils.ResponseJSON(w, http.StatusInternalServerError, "failed to generate token", nil)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    token,
		MaxAge:   15 * 60,
		Path:     "/",
		HttpOnly: true,
	})

	utils.ResponseJSON(w, http.StatusOK, "login successfully", map[string]string{
		"token": token,
	})

}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	// userID, ok := middlewares.GetUserID(r)
	// if !ok {
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	return
	// }

	emailID, ok := middlewares.GetUserEmail(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// json.NewEncoder(w).Encode(map[string]any{
	// 	"user-id":  userID,
	// 	"email-id": emailID,
	// })

	var user *models.UserResponse

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err := h.DB.Users.FindOne(ctx, bson.D{{Key: "email", Value: emailID}}).Decode(&user)
	if err != nil {
		utils.ResponseJSON(w, http.StatusOK, "couldn't find", nil)
		return
	}

	if err := utils.ResponseJSON(w, http.StatusOK, "user fetched", user); err != nil {
		log.Println("error encoding response")
	}

}
