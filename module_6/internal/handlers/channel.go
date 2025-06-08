package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"module_6/internal/models"
	"module_6/internal/services"
	"module_6/internal/utils"
)

type ChannelHandler struct {
	channelService *services.ChannelService
	wsService      *services.WebSocketService // Added for broadcasting new messages
}

func NewChannelHandler(s *services.ChannelService, wsS *services.WebSocketService) *ChannelHandler {
	return &ChannelHandler{
		channelService: s,
		wsService:      wsS,
	}
}

func (h *ChannelHandler) GetHistory(c *fiber.Ctx) error {
	log.Println("GetHistory handler called")
	req := new(models.GetHistoryRequest)
	if err := c.QueryParser(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters")
	}

	channelID, err := primitive.ObjectIDFromHex(req.ChannelID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Channel ID format")
	}

	// UserID from authenticated context (set by auth middleware)
	// userID := c.Locals("userID").(primitive.ObjectID) // Assuming userID is set by middleware

	messages, err := h.channelService.GetMessages(c.Context(), channelID, int64(req.Limit), int64(req.Offset))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve messages")
	}

	return utils.JSONResponse(c, fiber.StatusOK, "Messages retrieved successfully", messages)
}

func (h *ChannelHandler) SendMessage(c *fiber.Ctx) error {
	log.Println("SendMessage handler called")
	req := new(models.SendMessageRequest)
	if err := c.BodyParser(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	channelID, err := primitive.ObjectIDFromHex(req.ChannelID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Channel ID format")
	}

	// Get user ID from authenticated context
	userID, ok := c.Locals("userID").(primitive.ObjectID)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User ID not found in context")
	}
	username, ok := c.Locals("username").(string)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Username not found in context")
	}

	message, err := h.channelService.SendMessage(c.Context(), channelID, userID, username, req.Content)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to send message")
	}

	// Broadcast message to WebSocket clients (if wsService is available)
	if h.wsService != nil {
		h.wsService.BroadcastMessage(models.Message{
			ID:             message.ID,
			ChannelID:      message.ChannelID,
			SenderID:       message.SenderID,
			SenderUsername: message.SenderUsername,
			Content:        message.Content,
			Timestamp:      message.Timestamp,
		})
	}

	return utils.JSONResponse(c, fiber.StatusOK, "Message sent successfully", fiber.Map{
		"message_id": message.ID.Hex(),
		"timestamp":  message.Timestamp,
	})
}
