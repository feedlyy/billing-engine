package router

import (
	_const "billingg-engine/const"
	"billingg-engine/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

var jwtKey = []byte("very_secret_key")

// generateJWT generates a new JWT token for the given username
func generateJWT(username, role string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &model.UserClaims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// verifyJWT checks if the provided token is valid
func verifyJWT(tokenString string) (*model.UserClaims, error) {
	claims := &model.UserClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		log.Fatal(err)
	}

	return claims, nil
}

type Middleware struct {
	userDataSource []model.User
}

func NewMiddleware(data []model.User) Middleware {
	return Middleware{
		userDataSource: data,
	}
}

func (m Middleware) AuthenticateUser(username, password string) (string, error) {
	for _, user := range m.userDataSource {
		if user.Username == username {
			err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
			if err != nil {
				return "", fmt.Errorf("invalid credentials")
			}
			return generateJWT(username, user.Role)
		}
	}
	return "", fmt.Errorf("invalid credentials")
}

func (m Middleware) Login(w http.ResponseWriter, r *http.Request) {
	tokenString, err := m.AuthenticateUser(r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func AuthMiddlewareWithRole(next func(http.ResponseWriter, *http.Request), role string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing auth token", http.StatusUnauthorized)
			return
		}

		// Expecting token to start like "Bearer <token>"
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			http.Error(w, "Missing or invalid Bearer token", http.StatusUnauthorized)
			return
		}

		claims, err := verifyJWT(tokenStr)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		if claims.Role != role {
			http.Error(w, "you did not have permission to view this page", http.StatusForbidden)
			return
		}

		// Add user claims to the request context
		*r = *r.WithContext(context.WithValue(r.Context(), _const.UserContextKey, claims))

		// Token is valid, proceed to the next handler
		next(w, r)
	}
}
