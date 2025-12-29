package payments

import (
	"fmt"
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

func (h *Handler) StripeWebhook() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload struct {
			Type string `json:"type"`
			Data struct {
				Object struct {
					ID       string            `json:"id"`
					Metadata map[string]string `json:"metadata"`
				} `json:"object"`
			} `json:"data"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid payload"})
		}
		orderIDStr := payload.Data.Object.Metadata["order_id"]
		if orderIDStr == "" {
			return c.Status(200).JSON(fiber.Map{"message": "ignored"})
		}
		var oid uint
		if _, err := fmt.Sscan(orderIDStr, &oid); err != nil {
			return c.Status(200).JSON(fiber.Map{"message": "ignored"})
		}
		switch payload.Type {
		case "checkout.session.completed", "payment_intent.succeeded":
			_ = h.svc.UpdatePaymentStatus(c.Context(), oid, "paid")
		case "checkout.session.expired", "payment_intent.payment_failed":
			_ = h.svc.UpdatePaymentStatus(c.Context(), oid, "declined")
		default:
		}
		return c.SendStatus(200)
	}
}
