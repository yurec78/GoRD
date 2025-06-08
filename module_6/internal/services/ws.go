package services

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"module_6/internal/models"

	"github.com/gofiber/websocket/v2" // Corrected import
)

// WebSocketService manages WebSocket connections and message broadcasting
type WebSocketService struct {
	messagesCollection *mongo.Collection
	clients            map[*websocket.Conn]bool // Connected clients: conn -> true (dummy value)
	broadcast          chan models.Message      // Channel for broadcasting messages
	mu                 sync.Mutex               // Mutex for protecting access to the clients map
}

// NewWebSocketService creates a new WebSocketService
func NewWebSocketService(messagesCol *mongo.Collection) *WebSocketService {
	ws := &WebSocketService{
		messagesCollection: messagesCol,
		clients:            make(map[*websocket.Conn]bool),
		broadcast:          make(chan models.Message),
	}
	go ws.runBroadcaster() // Start the message broadcaster goroutine immediately
	return ws
}

// RegisterClient adds a new WebSocket connection to the pool
func (s *WebSocketService) RegisterClient(conn *websocket.Conn) {
	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()
	log.Printf("WebSocket client connected: %s. Total clients: %d", conn.RemoteAddr().String(), len(s.clients))
}

// UnregisterClient removes a WebSocket connection from the pool
func (s *WebSocketService) UnregisterClient(conn *websocket.Conn) {
	s.mu.Lock()
	delete(s.clients, conn)
	s.mu.Unlock()
	log.Printf("WebSocket client disconnected: %s. Total clients: %d", conn.RemoteAddr().String(), len(s.clients))
}

// BroadcastMessage sends a message to all connected clients (or specific channel clients, if implemented)
// This message is sent to an internal channel, which then gets processed by runBroadcaster.
func (s *WebSocketService) BroadcastMessage(msg models.Message) {
	s.broadcast <- msg
	log.Printf("Message put into broadcast queue: %s", msg.Content)
}

func (s *WebSocketService) runBroadcaster() { // 's' is the receiver here
	for {
		select {
		case message := <-s.broadcast:
			s.mu.Lock() // Lock to safely iterate and modify the clients map
			for client := range s.clients {
				// ПЕРЕВІРТЕ ЗМІНИ ТУТ: 'wsService' передається до горутини
				go func(wsService *WebSocketService, client *websocket.Conn, msg models.Message) {
					err := client.WriteJSON(msg) // msg теж передається
					if err != nil {
						log.Printf("Error writing to WebSocket client %s: %v. Unregistering client.", client.RemoteAddr().String(), err)
						client.Close()
						wsService.UnregisterClient(client) // <--- ТЕПЕР 'wsService' ВИЗНАЧЕНО!
					}
				}(s, client, message) // <--- ТУТ ПЕРЕДАЄМО 's' (яке стане 'wsService' всередині горутини), 'client' та 'message'
			}
			s.mu.Unlock()
		case <-time.After(10 * time.Second):
			// ... (rest of the code)
		}
	}
}

// ListenForNewMessagesFromDB (Placeholder for MongoDB Change Streams)
// This function would be started as a separate goroutine if you want real-time updates
// for messages inserted into the database by other services/applications.
func (s *WebSocketService) ListenForNewMessagesFromDB(ctx context.Context) {
	// ... (rest of the placeholder code remains the same)
	log.Println("WebSocketService: MongoDB Change Stream listener (placeholder) started.")
}
