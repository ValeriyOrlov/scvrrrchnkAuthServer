package handler

import (
	"errors"

	"github.com/ValeriyOrlov/scvrrrchnkAuthServer/internal/service"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	ctx := c.Context()
	user, err := h.authService.Register(ctx, req.Email, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidEmail) {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if errors.Is(err, service.ErrInvalidInput) {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if errors.Is(err, service.ErrShortUsername) {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if errors.Is(err, service.ErrWeakPassword) {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if errors.Is(err, service.ErrUserAlreadyExists) {
			return c.Status(409).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	ctx := c.Context()
	accessToken, refreshToken, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "bearer",
	})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad request"})
	}

	accessToken, refreshToken, err := h.authService.Refresh(c.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid refresh token"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "bearer",
	})
}
