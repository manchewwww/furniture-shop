package payments

import (
	"encoding/json"
	"fmt"
	"log"

	"furniture-shop/internal/config"
	payment_dto "furniture-shop/internal/dtos/payments"
	eo "furniture-shop/internal/entities/orders"
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

func (h *Handler) PayByCard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in payment_dto.CardPayment
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

		event, err := webhook.ConstructEvent(payload, sig, config.Env.StripeWebhookSecret)
		if err != nil {
			log.Printf("stripe webhook signature error: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid signature"})
		}

		update := func(orderID uint, paymentStatus, orderStatus string) {
			if err := h.svc.ProcessPaymentResult(c.Context(), orderID, paymentStatus, orderStatus); err != nil {
				log.Printf("ProcessPaymentResult failed (order=%d paid-status=%v order-status=%v): %v", orderID, paymentStatus, orderStatus, err)
			}
		}

		parseOID := func(s string) (uint, bool) {
			var oid uint
			if s == "" {
				return 0, false
			}
			if _, err := fmt.Sscan(s, &oid); err != nil {
				return 0, false
			}
			return oid, true
		}

		switch event.Type {

		case "payment_intent.succeeded":
			var pi stripe.PaymentIntent
			if err := json.Unmarshal(event.Data.Raw, &pi); err == nil {
				if oid, ok := parseOID(pi.Metadata["order_id"]); ok {
					update(oid, eo.PaymentStatusPaid, eo.OrderStatusProcessing)
				}
			}
		case "payment_intent.payment_failed":
			var pi stripe.PaymentIntent
			if err := json.Unmarshal(event.Data.Raw, &pi); err == nil {
				if oid, ok := parseOID(pi.Metadata["order_id"]); ok {
					update(oid, eo.PaymentStatusDeclined, eo.OrderStatusCancelled)
				}
			}
		case "checkout.session.expired":
			var sess stripe.CheckoutSession
			if err := json.Unmarshal(event.Data.Raw, &sess); err == nil {
				orderIDStr := sess.ClientReferenceID
				if orderIDStr == "" {
					orderIDStr = sess.Metadata["order_id"]
				}
				if oid, ok := parseOID(orderIDStr); ok {
					update(oid, eo.PaymentStatusCancelled, eo.OrderStatusCancelled)
				}
			}
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
