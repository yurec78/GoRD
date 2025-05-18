package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultMongoDBURI = "mongodb://localhost:27017"
	defaultDatabase   = "http_store_db"
	defaultServerPort = ":8080"
)

// Глобальна змінна для сховища, доступна для обробників
var store *MongoStore

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = defaultMongoDBURI
	}
	dbName := os.Getenv("MONGO_DATABASE")
	if dbName == "" {
		dbName = defaultDatabase
	}
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = defaultServerPort
	}
	if !strings.HasPrefix(serverPort, ":") { // Переконуємося, що порт починається з :
		serverPort = ":" + serverPort
	}

	// Створюємо контекст з тайм-аутом для підключення
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Підключення до MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Перевірка підключення
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v. Check connection and if MongoDB is running at %s.", err, mongoURI)
	}
	log.Printf("Successfully connected to MongoDB at %s", mongoURI)

	// Відкладене відключення клієнта
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		} else {
			log.Println("Disconnected from MongoDB.")
		}
	}()

	// Ініціалізація сховища
	store = NewMongoStore(client, dbName)
	log.Printf("Using database: %s", dbName)

	// Реєстрація обробників
	http.HandleFunc("/create_collection", makePostHandler(handleCreateCollection))
	http.HandleFunc("/list_collections", makePostHandler(handleListCollections))
	http.HandleFunc("/delete_collection", makePostHandler(handleDeleteCollection))

	http.HandleFunc("/put_document", makePostHandler(handlePutDocument))
	http.HandleFunc("/get_document", makePostHandler(handleGetDocument))
	http.HandleFunc("/list_documents", makePostHandler(handleListDocuments))
	http.HandleFunc("/delete_document", makePostHandler(handleDeleteDocument))

	http.HandleFunc("/create_index", makePostHandler(handleCreateIndex))
	http.HandleFunc("/delete_index", makePostHandler(handleDeleteIndex))

	// Запуск HTTP сервера
	log.Printf("Server starting on port %s", serverPort)
	if err := http.ListenAndServe(serverPort, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
