package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"module_6/internal/clients"
	"module_6/internal/config"
	"module_6/internal/handlers"
	"module_6/internal/middlewares"
	"module_6/internal/services"

	swagger "github.com/arsmn/fiber-swagger/v2"
	_ "module_6/docs"
)

func main() {
	// 1. Завантаження конфігурації
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 2. Ініціалізація MongoDB клієнта
	mongoClient, db, err := clients.InitMongoDB(cfg.MongoDBURI, cfg.MongoDBName)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer func() {
		if err = mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()
	log.Printf("Connected to MongoDB: %s/%s", cfg.MongoDBURI, cfg.MongoDBName)

	// 3. Ініціалізація сервісів
	authService := services.NewAuthService(db.Collection("users"), cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	channelService := services.NewChannelService(db.Collection("messages"), db.Collection("channels"))
	wsService := services.NewWebSocketService(db.Collection("messages")) // Для WebSocket

	// 4. Створення Fiber-додатку
	app := fiber.New()

	// 5. Додавання глобальних middleware
	app.Use(recover.New())
	app.Use(compress.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip} ${status} ${method} ${path} - ${latency}\n",
	}))
	app.Use(middlewares.RequestLogger())

	// 6. Налаштування Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// 7. Налаштування маршрутів (Endpoints)
	api := app.Group("/api/v1")

	// Ендпоінти авторизації
	authHandler := handlers.NewAuthHandler(authService)
	api.Post("/auth/sign-up", authHandler.SignUp)
	api.Post("/auth/sign-in", authHandler.SignIn)
	api.Post("/auth/refresh-token", authHandler.RefreshToken)

	// Ендпоінти чату (захищені авторизацією)
	channelHandler := handlers.NewChannelHandler(channelService, wsService)
	authMiddleware := middlewares.NewAuthMiddleware(cfg.JWTSecret)
	api.Get("/channel/history", authMiddleware, channelHandler.GetHistory)
	api.Post("/channel/send", authMiddleware, channelHandler.SendMessage)

	// WebSocket ендпоінт
	wsHandler := handlers.NewWebSocketHandler(wsService, authService, cfg.JWTSecret)
	app.Use("/channel/listen", wsHandler.UpgradeWebSocket)

	// 8. Запуск сервера
	listenAddr := ":" + cfg.Port
	log.Printf("Starting Fiber server on %s", listenAddr)

	// Запуск сервера в горутині, щоб мати змогу коректно вимкнути його
	go func() {
		if err := app.Listen(listenAddr); err != nil && err.Error() != "http: Server closed" {
			log.Fatalf("Fiber server error: %v", err)
		}
	}()

	// 9. Коректне вимкнення сервера
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down Fiber server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Fiber server shutdown error: %v", err)
	}
	log.Println("Fiber server gracefully stopped.")
}
