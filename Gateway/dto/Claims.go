package dto

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	Role     string `json:"role"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
