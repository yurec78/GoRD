package handlers_test // Окремий тестовий пакет

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"module_6/internal/clients"
	"module_6/internal/config"
	"module_6/internal/handlers"
	"module_6/internal/models"
	"module_6/internal/services" // Імпортуємо реальний сервіс
	"net/http/httptest"
	"strings"
	"testing"
)

// SetupTestApp створює Fiber-додаток для тестування
func SetupTestApp(t *testing.T, mongoClient *mongo.Client, cfg *config.Config) *fiber.App {
	app := fiber.New()

	db := mongoClient.Database(cfg.MongoDBName)
	usersCol := db.Collection("users")
	// Очистити колекцію користувачів перед тестом, щоб забезпечити чисте середовище
	if _, err := usersCol.DeleteMany(context.Background(), bson.M{}); err != nil {
		t.Fatalf("Failed to clear users collection: %v", err)
	}

	authService := services.NewAuthService(usersCol, cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	authHandler := handlers.NewAuthHandler(authService)

	app.Post("/auth/sign-up", authHandler.SignUp)
	app.Post("/auth/sign-in", authHandler.SignIn)

	return app
}

func TestAuthHandler_SignUp(t *testing.T) {
	// Load config for DB connection
	cfg, err := config.LoadConfig()
	assert.NoError(t, err)

	// Connect to test MongoDB
	mongoClient, _, err := clients.InitMongoDB(cfg.MongoDBURI, cfg.MongoDBName)
	assert.NoError(t, err)
	defer mongoClient.Disconnect(context.Background())

	app := SetupTestApp(t, mongoClient, cfg)

	tests := []struct {
		name         string
		body         models.SignUpRequest
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "Successful sign up",
			body:         models.SignUpRequest{Username: "testuser1", Password: "password123"},
			expectedCode: fiber.StatusOK,
			expectedMsg:  "User registered successfully",
		},
		{
			name:         "Sign up with existing username",
			body:         models.SignUpRequest{Username: "testuser1", Password: "password123"}, // Use existing user
			expectedCode: fiber.StatusInternalServerError,                                      // Or Bad Request if your service returns specific errors
			expectedMsg:  "Failed to register user",                                            // Or specific error from service
		},
		{
			name:         "Invalid request body (missing username)",
			body:         models.SignUpRequest{Password: "password123"},
			expectedCode: fiber.StatusBadRequest,
			expectedMsg:  "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/auth/sign-up", strings.NewReader(string(requestBody)))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1) // -1 for no timeout
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			var response map[string]string
			json.NewDecoder(resp.Body).Decode(&response)
			assert.Contains(t, response["message"], tt.expectedMsg)
		})
	}
}

func TestAuthHandler_SignIn(t *testing.T) {
	// Setup (similar to SignUp, ensure a user exists)
	cfg, err := config.LoadConfig()
	assert.NoError(t, err)

	mongoClient, _, err := clients.InitMongoDB(cfg.MongoDBURI, cfg.MongoDBName)
	assert.NoError(t, err)
	defer mongoClient.Disconnect(context.Background())

	app := SetupTestApp(t, mongoClient, cfg)

	// Register a user first for sign-in test
	authService := services.NewAuthService(mongoClient.Database(cfg.MongoDBName).Collection("users"), cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	err = authService.RegisterUser(context.Background(), "signinuser", "signinpassword")
	assert.NoError(t, err)

	tests := []struct {
		name         string
		body         models.SignInRequest
		expectedCode int
		expectedMsg  string // For error messages
		expectToken  bool
	}{
		{
			name:         "Successful sign in",
			body:         models.SignInRequest{Username: "signinuser", Password: "signinpassword"},
			expectedCode: fiber.StatusOK,
			expectedMsg:  "Login successful",
			expectToken:  true,
		},
		{
			name:         "Invalid credentials (wrong password)",
			body:         models.SignInRequest{Username: "signinuser", Password: "wrongpassword"},
			expectedCode: fiber.StatusUnauthorized, // Or specific error from service
			expectedMsg:  "Invalid credentials",
			expectToken:  false,
		},
		// ... more test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/auth/sign-in", strings.NewReader(string(requestBody)))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			var authResp models.AuthResponse
			if tt.expectToken {
				json.NewDecoder(resp.Body).Decode(&authResp)
				assert.NotEmpty(t, authResp.Token)
				assert.NotEmpty(t, authResp.RefreshToken)
			} else {
				var errorResp map[string]string
				json.NewDecoder(resp.Body).Decode(&errorResp)
				assert.Contains(t, errorResp["error"], tt.expectedMsg)
				assert.Empty(t, authResp.Token)
			}
		})
	}
}
