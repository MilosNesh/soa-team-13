package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"stakeholders.com/dto"
)

var jwtSecret = []byte("sekret_key_12#4") // uzmi iz env-a u realnom kodu

// ExtractUser vrati (userId, username, role) iz Authorization headera.
func ExtractUser(r *http.Request) (string, string, string, error) {
	authz := r.Header.Get("Authorization")
	if !strings.HasPrefix(authz, "Bearer ") {
		return "", "", "", errors.New("missing bearer token")
	}
	tokenStr := strings.TrimPrefix(authz, "Bearer ")

	claims := &dto.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", "", "", errors.New("invalid token")
	}
	if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		return "", "", "", errors.New("token expired")
	}

	return claims.Subject, claims.Username, claims.Role, nil
}
