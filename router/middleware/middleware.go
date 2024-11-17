package middleware

import (
	_const "billingg-engine/const"
	"billingg-engine/model"
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

var jwtKey = []byte("very_secret_key")

func GenerateJWT(username, role string) (string, error) {
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

func verifyJWT(tokenString string) (*model.UserClaims, error) {
	claims := &model.UserClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func AuthMiddlewareWithRole(next func(http.ResponseWriter, *http.Request), role string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")

		type Err struct {
			Error string `json:"error"`
		}

		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Err{Error: "Missing Auth Token"})
			return
		}

		// Expecting token to start like "Bearer <token>"
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Err{Error: "Missing or invalid Bearer token"})
			return
		}

		claims, err := verifyJWT(tokenStr)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Err{Error: "Invalid or expired token"})
			return
		}

		if claims.Role != role {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Err{Error: "you did not have permission to view this page"})
			return
		}

		// Add user claims to the request context
		*r = *r.WithContext(context.WithValue(r.Context(), _const.UserContextKey, claims))

		// Token is valid, proceed to the next handler
		next(w, r)
	}
}
