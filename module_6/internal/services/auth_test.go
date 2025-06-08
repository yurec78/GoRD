package services_test // Використовуємо _test суфікс для окремого тестового пакета

import (
	"context"
	"errors"
	"testing"
	"time"

	"module_6/internal/config" // Якщо вам потрібен доступ до конфігурації
	"module_6/internal/models"
	"module_6/internal/services" // Імпортуємо пакет, який тестуємо

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// Mokery для генерації моків
	"github.com/stretchr/testify/assert" // Популярна бібліотека для асертів
	"github.com/stretchr/testify/mock"   // Для моків
)

// --- Mocking MongoDB Collection ---
// Це приклад того, як можна створити "мок" для mongo.Collection
// Для складніших сценаріїв можна використовувати бібліотеки типу 'gomock'
// або 'testify/mock' з автоматичною генерацією.
// Тут ми робимо ручний мок для простоти.

type MockMongoCollection struct {
	mock.Mock
}

// Implement required methods for mongo.Collection that AuthService uses
func (m *MockMongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockMongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.SingleResult)
}

// MockSingleResult for FindOne
type MockSingleResult struct {
	mock.Mock
	err error
	doc interface{}
}

func (m *MockSingleResult) Decode(val interface{}) error {
	if m.err != nil {
		return m.err
	}
	// Simulate decoding into the passed value
	if m.doc != nil {
		docBytes, _ := bson.Marshal(m.doc)
		bson.Unmarshal(docBytes, val)
	}
	return nil
}

func (m *MockSingleResult) Err() error {
	return m.err
}

// --- Test Suite for AuthService ---

func TestAuthService_RegisterUser(t *testing.T) {
	// 1. Setup - Створення моків та сервісу
	mockUsersCol := new(MockMongoCollection)
	cfg := &config.Config{
		JWTSecret:       "testsecret",
		AccessTokenTTL:  3600,
		RefreshTokenTTL: 2592000,
	}
	authService := services.NewAuthService(mockUsersCol, cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)

	// 2. Define Test Cases - Визначення тестових сценаріїв
	tests := []struct {
		name        string
		username    string
		password    string
		mockSetup   func()
		expectedErr error
	}{
		{
			name:     "Successful registration",
			username: "testuser",
			password: "password123",
			mockSetup: func() {
				mockUsersCol.On("InsertOne", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("models.User")).Return(&mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:     "Registration with existing username",
			username: "existinguser",
			password: "password123",
			mockSetup: func() {
				// Simulate duplicate key error
				mockUsersCol.On("InsertOne", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("models.User")).Return(nil, mongo.WriteException{WriteErrors: []mongo.WriteError{{Code: 11000}}}).Once()
			},
			expectedErr: errors.New("username already exists"), // Перевірте, яку саме помилку повертає ваш сервіс
		},
		{
			name:     "Registration with DB error",
			username: "dbuser",
			password: "password123",
			mockSetup: func() {
				mockUsersCol.On("InsertOne", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("models.User")).Return(nil, errors.New("database error")).Once()
			},
			expectedErr: errors.New("database error"),
		},
	}

	// 3. Run Tests - Запуск тестів
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock calls for each test case
			mockUsersCol.Calls = []mock.Call{}
			tt.mockSetup() // Налаштування моків для поточного тестового сценарію

			err := authService.RegisterUser(context.Background(), tt.username, tt.password)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error()) // Перевірка, що помилка містить очікуване повідомлення
			} else {
				assert.NoError(t, err)
			}
			mockUsersCol.AssertExpectations(t) // Перевірка, що моки були викликані як очікувалося
		})
	}
}

func TestAuthService_AuthenticateUser(t *testing.T) {
	// Setup
	mockUsersCol := new(MockMongoCollection)
	cfg := &config.Config{
		JWTSecret:       "testsecret",
		AccessTokenTTL:  3600,
		RefreshTokenTTL: 2592000,
	}
	authService := services.NewAuthService(mockUsersCol, cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)

	// Test Cases
	tests := []struct {
		name                 string
		username             string
		password             string
		mockSetup            func()
		expectedToken        string
		expectedRefreshToken string
		expectedErr          error
	}{
		{
			name:     "Successful authentication",
			username: "validuser",
			password: "validpassword",
			mockSetup: func() {
				hashedPassword, _ := services.HashPassword("validpassword") // Хешуємо пароль
				mockUser := models.User{
					ID: primitive.NewObjectID(), Username: "validuser", Password: hashedPassword, CreatedAt: time.Now(), UpdatedAt: time.Now(),
				}
				mockResult := &MockSingleResult{doc: mockUser}
				mockUsersCol.On("FindOne", mock.Anything, bson.M{"username": "validuser"}).Return(mockResult).Once()
			},
			expectedErr: nil,
		},
		{
			name:     "User not found",
			username: "nonexistent",
			password: "anypassword",
			mockSetup: func() {
				mockResult := &MockSingleResult{err: mongo.ErrNoDocuments}
				mockUsersCol.On("FindOne", mock.Anything, bson.M{"username": "nonexistent"}).Return(mockResult).Once()
			},
			expectedErr: errors.New("user not found"),
		},
		// Додайте більше тест-кейсів: неправильний пароль, помилка БД при пошуку, тощо
	}

	// Run Tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsersCol.Calls = []mock.Call{} // Clear mocks
			tt.mockSetup()

			token, refreshToken, err := authService.AuthenticateUser(context.Background(), tt.username, tt.password)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
				assert.Empty(t, token)
				assert.Empty(t, refreshToken)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.NotEmpty(t, refreshToken)
			}
			mockUsersCol.AssertExpectations(t)
		})
	}
}

// Helper function that might be in services/auth.go (example)
// This should be part of your actual service logic, not just for tests.
func TestAuthService_HashPassword(t *testing.T) {
	password := "testpassword"
	hashedPassword, err := services.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	assert.True(t, services.CheckPasswordHash(password, hashedPassword))
	assert.False(t, services.CheckPasswordHash("wrongpassword", hashedPassword))
}
