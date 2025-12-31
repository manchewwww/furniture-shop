package orders

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	stripe "github.com/stripe/stripe-go/v84"
	session "github.com/stripe/stripe-go/v84/checkout/session"

	"furniture-shop/internal/config"
	order_dto "furniture-shop/internal/dtos/orders"
	"furniture-shop/internal/entities/orders"
	"furniture-shop/internal/service"
	"furniture-shop/internal/service/mailer"
	vld "furniture-shop/internal/validation"
)

type Handler struct {
	svc service.OrdersService
}

func NewOrdersHandler(svc service.OrdersService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateOrder() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in order_dto.CreateOrderInput
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
		}
		if err := vld.ValidateStruct(in); err != nil {
			return err
		}

		if uid, ok := c.Locals("user_id").(uint); ok {
			in.UserID = &uid
		}

		order, err := h.svc.CreateOrder(c.Context(), in)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"message": err.Error()})
		}

		to := c.Locals("user_email")
		if s, ok := to.(string); ok && s != "" {
			mailer.NewSender().Send(s, "Order created", fmt.Sprintf("Your order #%d has been created and is pending.", order.ID))
		} else if in.Email != "" {
			mailer.NewSender().Send(in.Email, "Order created", fmt.Sprintf("Your order #%d has been created and is pending.", order.ID))
		}

		if in.PaymentMethod == "card" {
			url, err := h.createStripeSession(order)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"message": "payment init failed"})
			}
			return c.JSON(fiber.Map{
				"order_id":                       order.ID,
				"checkout_url":                   url,
				"estimated_production_time_days": order.EstimatedProductionTimeDays,
			})
		}

		return c.JSON(order)
	}
}

func (h *Handler) PayExistingOrder() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, _ := c.Locals("user_id").(uint)
		id, err := h.getID(c)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid id"})
		}

		order, err := h.svc.GetUserOrder(c.Context(), uid, id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"message": "order not found"})
		}

		url, err := h.createStripeSession(order)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "stripe error"})
		}

		return c.JSON(fiber.Map{"checkout_url": url})
	}
}

func (h *Handler) UserOrders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := c.Locals("user_id").(uint)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"message": "unauthorized"})
		}
		orders, err := h.svc.ListUserOrders(c.Context(), uid)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(orders)
	}
}

func (h *Handler) UserOrderDetails() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := c.Locals("user_id").(uint)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"message": "unauthorized"})
		}
		var id uint
		if _, err := fmt.Sscan(c.Params("id"), &id); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid id"})
		}
		order, err := h.svc.GetUserOrder(c.Context(), uid, id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"message": "not found"})
		}
		return c.JSON(order)
	}
}

func (h *Handler) AdminListOrders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		status := c.Query("status")
		orders, err := h.svc.AdminListOrders(c.Context(), status)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(orders)
	}
}

func (h *Handler) AdminUpdateOrderStatus() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in order_dto.PatchStatusRequest
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
		}
		var id uint
		if _, err := fmt.Sscan(c.Params("id"), &id); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid id"})
		}
		if err := h.svc.AdminUpdateOrderStatus(c.Context(), id, in.Status); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "updated"})
	}
}

func (h *Handler) getID(c *fiber.Ctx) (uint, error) {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func (h *Handler) createStripeSession(order *orders.Order) (string, error) {
	stripe.Key = config.Env.StripeSecretKey
	fe := "http://localhost:5173"

	amount := int64(math.Round(order.TotalPrice * 100))
	expiresAt := time.Now().Add(30 * time.Minute).Unix()
	orderIDStr := strconv.Itoa(int(order.ID))

	params := &stripe.CheckoutSessionParams{
		Mode:              stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:        stripe.String(fmt.Sprintf("%s/payment/success?session_id={CHECKOUT_SESSION_ID}&order_id=%d", fe, order.ID)),
		CancelURL:         stripe.String(fmt.Sprintf("%s/payment/cancel?order_id=%d", fe, order.ID)),
		ClientReferenceID: stripe.String(orderIDStr),
		ExpiresAt:         stripe.Int64(expiresAt),
		Metadata:          map[string]string{"order_id": orderIDStr},
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: map[string]string{"order_id": orderIDStr},
		},
		LineItems: []*stripe.CheckoutSessionLineItemParams{{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("eur"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(fmt.Sprintf("Order #%d", order.ID)),
				},
				UnitAmount: stripe.Int64(amount),
			},
			Quantity: stripe.Int64(1),
		}},
	}

	sess, err := session.New(params)
	if err != nil {
		return "", err
	}
	return sess.URL, nil
}
