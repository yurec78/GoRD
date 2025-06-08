package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port            string
	MongoDBURI      string
	MongoDBName     string
	JWTSecret       string
	AccessTokenTTL  int // в секундах
	RefreshTokenTTL int // в секундах
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		MongoDBURI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDBName:     getEnv("MONGODB_NAME", "chatdb"),
		JWTSecret:       getEnv("JWT_SECRET", "supersecretjwtkey"), // Дуже важливо змінити для продакшену
		AccessTokenTTL:  getEnvAsInt("ACCESS_TOKEN_TTL", 3600),     // 1 година
		RefreshTokenTTL: getEnvAsInt("REFRESH_TOKEN_TTL", 2592000), // 30 днів
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set. This is critical for security.")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	strValue := getEnv(key, "")
	if strValue == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		log.Printf("Warning: Environment variable %s is not a valid integer, using default %d. Error: %v", key, defaultValue, err)
		return defaultValue
	}
	return intValue
}
