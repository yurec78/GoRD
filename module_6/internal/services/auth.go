package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt" // Для хешування паролів
	"module_6/internal/models"
	"module_6/internal/utils" // Для JWT
)

// AuthService handles user authentication logic
type AuthService struct {
	usersCollection *mongo.Collection
	jwtSecret       string
	accessTokenTTL  int
	refreshTokenTTL int
}

// NewAuthService creates a new AuthService
func NewAuthService(usersCol *mongo.Collection, secret string, accessTTL, refreshTTL int) *AuthService {
	return &AuthService{
		usersCollection: usersCol,
		jwtSecret:       secret,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

// RegisterUser handles the user registration business logic
func (s *AuthService) RegisterUser(ctx context.Context, username, password string) error {
	// Перевірка, чи користувач вже існує
	count, err := s.usersCollection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		log.Printf("Error checking for existing user: %v", err)
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if count > 0 {
		return errors.New("username already exists")
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := models.User{
		ID:        primitive.NewObjectID(),
		Username:  username,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.usersCollection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		// Перевірка на дублювання ключа (якщо MongoDB повертає Code 11000)
		if writeException, ok := err.(mongo.WriteException); ok {
			for _, writeError := range writeException.WriteErrors {
				if writeError.Code == 11000 { // Duplicate key error
					return errors.New("username already exists")
				}
			}
		}
		return fmt.Errorf("failed to register user: %w", err)
	}
	log.Printf("User registered: %s", username)
	return nil
}

// AuthenticateUser handles user authentication business logic
func (s *AuthService) AuthenticateUser(ctx context.Context, username, password string) (string, string, error) {
	var user models.User
	err := s.usersCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", "", errors.New("user not found")
		}
		log.Printf("Error finding user: %v", err)
		return "", "", fmt.Errorf("failed to find user: %w", err)
	}

	if !CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("invalid password")
	}

	// Генерація токенів
	accessToken, err := utils.GenerateJWT(user.ID, user.Username, s.jwtSecret, s.accessTokenTTL)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, err := utils.GenerateRefreshToken() // GenerateRefreshToken in utils/jwt.go
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	log.Printf("User authenticated: %s", username)
	return accessToken, refreshToken, nil
}

// RefreshTokens handles refreshing of access and refresh tokens
func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	// In a real application, you would validate the refresh token against a database.
	// For simplicity, we'll just check if it's not empty.
	if refreshToken == "" {
		return "", "", errors.New("refresh token missing")
	}

	// Assuming a simplified lookup or re-issuance logic for testing
	// In a real scenario, you'd find the user associated with this refresh token
	// and issue new tokens.
	// For now, let's just create new dummy tokens if the refresh token is valid.
	// This would require a more robust refresh token storage/validation strategy.
	log.Printf("Refreshing tokens for refresh token: %s", refreshToken)

	// Dummy user ID and username for demonstration, replace with actual user lookup
	dummyUserID := primitive.NewObjectID()
	dummyUsername := "refreshed_user"

	accessToken, err := utils.GenerateJWT(dummyUserID, dummyUsername, s.jwtSecret, s.accessTokenTTL)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new access token: %w", err)
	}
	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	return accessToken, newRefreshToken, nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
