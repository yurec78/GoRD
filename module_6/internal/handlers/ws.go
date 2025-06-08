package handlers

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"module_6/internal/models"
	"module_6/internal/services"
	"module_6/internal/utils"
)

type WebSocketHandler struct {
	wsService   *services.WebSocketService
	authService *services.AuthService
	jwtSecret   string
}

func NewWebSocketHandler(wsS *services.WebSocketService, authS *services.AuthService, jwtSecret string) *WebSocketHandler {
	return &WebSocketHandler{
		wsService:   wsS,
		authService: authS,
		jwtSecret:   jwtSecret,
	}
}

func (h *WebSocketHandler) UpgradeWebSocket(c *fiber.Ctx) error {

	token := c.Get("Authorization")
	if token == "" {
		token = c.Query("token")
	}

	if token == "" {
		log.Println("WebSocket upgrade: Authorization token is missing")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authorization token is missing")
	}
	if len(token) > 7 && strings.ToLower(token[:7]) == "bearer " {
		token = token[7:]
	}

	claims, err := utils.ParseJWT(token, h.jwtSecret) // Use your JWT utility
	if err != nil {
		log.Printf("WebSocket auth failed: %v", err)
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token")
	}

	c.Locals("userID", claims.UserID)
	c.Locals("username", claims.Username)

	if websocket.IsWebSocketUpgrade(c) {

		log.Println("WebSocket Upgrade requested and authenticated.")

		return websocket.New(func(conn *websocket.Conn) {

			userID, ok := c.Locals("userID").(string) // Claims.UserID is string
			if !ok {
				log.Println("Error: User ID not found in WebSocket context after auth")
				conn.Close()
				return
			}
			username, ok := c.Locals("username").(string) // Claims.Username is string
			if !ok {
				log.Println("Error: Username not found in WebSocket context after auth")
				conn.Close()
				return
			}
			objectUserID, _ := primitive.ObjectIDFromHex(userID)

			h.handleWebSocketConnection(conn, objectUserID, username)
		})(c)
	}
	return fiber.ErrUpgradeRequired
}

func (h *WebSocketHandler) handleWebSocketConnection(conn *websocket.Conn, userID primitive.ObjectID, username string) {
	log.Printf("WebSocket connection established for user: %s (%s)", username, userID.Hex())
	h.wsService.RegisterClient(conn)
	defer func() {
		h.wsService.UnregisterClient(conn)
		conn.Close()
	}()

	for {
		messageType, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error for user %s (%s): %v", username, userID.Hex(), err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("Client %s closed connection unexpectedly", username)
			}
			break
		}

		if messageType == websocket.TextMessage {
			log.Printf("Received WS message from %s: %s", username, string(msgBytes))

			var wsMsg models.WebSocketMessage
			if err := conn.ReadJSON(&wsMsg); err != nil {
				log.Printf("Error parsing WS message from %s: %v", username, err)
				continue // Skip to next message
			}

			channelObjectID, err := primitive.ObjectIDFromHex(wsMsg.ChannelID)
			if err != nil {
				log.Printf("Error converting channel ID '%s' to ObjectID: %v", wsMsg.ChannelID, err)
				continue
			}

			h.wsService.BroadcastMessage(models.Message{
				ChannelID:      channelObjectID,
				SenderID:       userID,
				SenderUsername: username,
				Content:        wsMsg.Content,
				Timestamp:      time.Now(), // Or savedMessage.Timestamp if from DB
			})

		}
	}
}
