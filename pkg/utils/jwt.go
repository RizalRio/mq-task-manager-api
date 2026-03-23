package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GenerateToken membuat JWT yang berisi ID user dan Role (future-proof untuk fitur Admin)
func GenerateToken(userID uuid.UUID, role string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")

	// Payload data yang disimpan di dalam token (Claims)
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		// Token akan kedaluwarsa dalam 24 jam (best practice keamanan)
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}