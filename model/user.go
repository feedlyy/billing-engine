package model

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Username     string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string
	Role         string `json:"role"`
	Tz
}

type UserClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
