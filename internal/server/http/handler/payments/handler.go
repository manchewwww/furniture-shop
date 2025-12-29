package payments

import (
	"furniture-shop/internal/service"
	vld "furniture-shop/internal/validation"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	svc service.PaymentService
}

func NewPaymentsHandler(svc service.PaymentService) *Handler {
	return &Handler{svc: svc}
}

type cardDTO = service.CardPayment

func (h *Handler) PayByCard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in cardDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
		}
		if err := vld.ValidateStruct(in); err != nil {
			return err
		}
		status, err := h.svc.PayByCard(c.Context(), in)
		if err != nil {
			return c.Status(402).JSON(fiber.Map{"message": "payment failed", "payment_status": "declined"})
		}
		return c.JSON(fiber.Map{"message": "payment accepted", "payment_status": status})
	}
}
