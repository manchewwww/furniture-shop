package payments

import (
	"encoding/json"
	"fmt"

	"furniture-shop/internal/config"
	"furniture-shop/internal/service"
	vld "furniture-shop/internal/validation"

	"github.com/gofiber/fiber/v2"
	stripe "github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/webhook"
)

type Handler struct {
	svc    service.PaymentService
	orders service.OrdersService
}

func NewPaymentsHandler(svc service.PaymentService, orders service.OrdersService) *Handler {
	return &Handler{svc: svc, orders: orders}
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
		payload := c.BodyRaw()
		sig := c.Get("Stripe-Signature")
		evt, err := webhook.ConstructEvent(payload, sig, config.Env.StripeWebhookSecret)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid signature"})
		}

		var orderIDStr string
		switch evt.Type {
		case "checkout.session.completed", "checkout.session.async_payment_succeeded", "checkout.session.async_payment_failed", "checkout.session.expired":
			var sess stripe.CheckoutSession
			if err := json.Unmarshal(evt.Data.Raw, &sess); err == nil {
				orderIDStr = sess.Metadata["order_id"]
			}
		case "payment_intent.succeeded", "payment_intent.payment_failed":
			var pi stripe.PaymentIntent
			if err := json.Unmarshal(evt.Data.Raw, &pi); err == nil {
				orderIDStr = pi.Metadata["order_id"]
			}
		default:
			return c.SendStatus(200)
		}

		if orderIDStr == "" {
			return c.SendStatus(200)
		}
		var oid uint
		if _, err := fmt.Sscan(orderIDStr, &oid); err != nil {
			return c.SendStatus(200)
		}

		switch evt.Type {
		case "checkout.session.completed", "payment_intent.succeeded", "checkout.session.async_payment_succeeded":
			_ = h.svc.UpdatePaymentStatus(c.Context(), oid, "paid")
			_ = h.orders.AdminUpdateOrderStatus(c.Context(), oid, "processing")
		case "payment_intent.payment_failed", "checkout.session.async_payment_failed", "checkout.session.expired":
			_ = h.svc.UpdatePaymentStatus(c.Context(), oid, "declined")
			_ = h.orders.AdminUpdateOrderStatus(c.Context(), oid, "cancelled")
		}
		return c.SendStatus(200)
	}
}
