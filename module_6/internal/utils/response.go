package utils

import "github.com/gofiber/fiber/v2"

// JSONResponse повертає стандартизовану JSON відповідь
func JSONResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"status":  statusCode,
		"message": message,
		"data":    data,
	})
}

// ErrorResponse повертає стандартизовану JSON відповідь з помилкою
func ErrorResponse(c *fiber.Ctx, statusCode int, errMessage string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"status":  statusCode,
		"message": "Error",
		"error":   errMessage,
	})
}
