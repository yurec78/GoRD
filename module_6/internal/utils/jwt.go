package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

// GenerateJWT generates a new JWT token
func GenerateJWT(userID primitive.ObjectID, username, secret string, ttl int) (string, error) {
	expirationTime := time.Now().Add(time.Duration(ttl) * time.Second)
	claims := &Claims{
		UserID:   userID.Hex(),
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

// ParseJWT parses and validates a JWT token
func ParseJWT(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

// GenerateRefreshToken generates a new refresh token (can be a simple UUID or a longer JWT)
func GenerateRefreshToken() (string, error) {
	// For simplicity, just return a new UUID. In a real app, you'd store it in DB.
	return primitive.NewObjectID().Hex(), nil // Using ObjectID as a simple unique string
}
