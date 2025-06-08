package handlers

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"module_6/internal/models"
	"module_6/internal/services"
	"module_6/internal/utils"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: s}
}

func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	log.Println("SignUp handler called") // DEBUG
	req := new(models.SignUpRequest)
	if err := c.BodyParser(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	return utils.JSONResponse(c, fiber.StatusOK, "User registered successfully", nil)
}

func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
	log.Println("SignIn handler called") // DEBUG
	req := new(models.SignInRequest)
	if err := c.BodyParser(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	return utils.JSONResponse(c, fiber.StatusOK, "Login successful", fiber.Map{"token": "dummy_token", "refresh_token": "dummy_refresh_token"})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	log.Println("RefreshToken handler called") // DEBUG
	req := new(models.AuthResponse)
	if err := c.BodyParser(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if req.RefreshToken == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Refresh token is missing")
	}

	return utils.JSONResponse(c, fiber.StatusOK, "Token refreshed successfully", fiber.Map{
		"token":         "new_dummy_token",
		"refresh_token": "new_dummy_refresh_token",
	})
}
