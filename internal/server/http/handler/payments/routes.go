package payments

import "github.com/gofiber/fiber/v2"

func Register(api fiber.Router, h *Handler) {
	api.Post("/webhooks/stripe", h.StripeWebhook())
}
