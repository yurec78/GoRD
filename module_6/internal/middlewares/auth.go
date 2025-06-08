package middlewares

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"module_6/internal/utils"
)

func NewAuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authorization header missing")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authorization header format must be Bearer {token}")
		}

		token := parts[1]
		claims, err := utils.ParseJWT(token, jwtSecret)
		if err != nil {
			log.Printf("JWT parsing error: %v", err)
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}

		// Store user ID and username in context locals for later use
		userID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid UserID in token claims")
		}
		c.Locals("userID", userID)
		c.Locals("username", claims.Username)

		return c.Next()
	}
}
