package middleware

import (
	"billingg-engine/model"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

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
			return GenerateJWT(username, user.Role)
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
