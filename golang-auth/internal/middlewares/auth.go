package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/abrarr21/golang-auth/internal/utils"
)

type contextKey string

const (
	UserKey      contextKey = "UserID"
	UserEmailKey contextKey = "userEmail"
)

func RequireAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			tokenString := extractToken(r)

			if tokenString == "" {
				utils.ResponseJSON(w, http.StatusUnauthorized, "missing token", nil)
				return
			}

			claims, err := utils.ParseToken(tokenString, secret)
			if err != nil {
				if errors.Is(err, utils.ErrTokenExpired) {
					utils.ResponseJSON(w, http.StatusUnauthorized, "token has expired", nil)
					return
				}
				utils.ResponseJSON(w, http.StatusUnauthorized, "invalid token", nil)
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func extractToken(r *http.Request) string {
	// For debugging
	// log.Println("Cookie header: ", r.Header.Get("Cookie"))
	// log.Println("Authorization header: ", r.Header.Get("Authorization"))

	cookie, err := r.Cookie("accessToken")
	if err == nil {
		return cookie.Value
	}

	_, token, found := strings.Cut(r.Header.Get("Authorization"), "Bearer ")
	if found {
		return token
	}

	return ""
}

func GetUserID(r *http.Request) (string, bool) {
	v, ok := r.Context().Value(UserKey).(string)
	return v, ok
}

func GetUserEmail(r *http.Request) (string, bool) {
	v, ok := r.Context().Value(UserEmailKey).(string)

	return v, ok
}
