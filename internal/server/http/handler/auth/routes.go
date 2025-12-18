package auth

import "github.com/gofiber/fiber/v2"

func Register(api fiber.Router, h *Handler) {
	api.Post("/auth/register", h.Register())
	api.Post("/auth/login", h.Login())
}
