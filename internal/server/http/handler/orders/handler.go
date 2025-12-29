package orders

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	stripe "github.com/stripe/stripe-go/v84"
	session "github.com/stripe/stripe-go/v84/checkout/session"

	"furniture-shop/internal/config"
	"furniture-shop/internal/service"
	vld "furniture-shop/internal/validation"
)

type Handler struct {
	svc service.OrdersService
}

func NewOrdersHandler(svc service.OrdersService) *Handler {
	return &Handler{svc: svc}
}

type orderItemIn struct {
	ProductID uint                     `json:"product_id" validate:"required,gt=0"`
	Quantity  int                      `json:"quantity" validate:"required,gt=0"`
	Options   []service.SelectedOption `json:"options"`
}
type createOrderDTO struct {
	UserID        *uint         `json:"user_id"`
	Name          string        `json:"name" validate:"required,min=2"`
	Email         string        `json:"email" validate:"required,email"`
	Address       string        `json:"address" validate:"required,min=5"`
	Phone         string        `json:"phone" validate:"required,phone"`
	Items         []orderItemIn `json:"items" validate:"required,min=1,dive"`
	PaymentMethod string        `json:"payment_method" validate:"required,oneof=card bank"`
}

func (h *Handler) CreateOrder() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in createOrderDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
		}
		if err := vld.ValidateStruct(in); err != nil {
			return err
		}
		if uidVal := c.Locals("user_id"); uidVal != nil {
			if uid, ok := uidVal.(uint); ok {
				in.UserID = &uid
			}
		}
		items := make([]service.CreateOrderItem, 0, len(in.Items))
		for _, it := range in.Items {
			items = append(items, service.CreateOrderItem{ProductID: it.ProductID, Quantity: it.Quantity, Options: it.Options})
		}
		order, err := h.svc.CreateOrder(c.Context(), service.CreateOrderInput{
			UserID: in.UserID, Name: in.Name, Email: in.Email, Address: in.Address, Phone: in.Phone,
			Items: items, PaymentMethod: in.PaymentMethod,
		})
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"message": err.Error()})
		}
		switch in.PaymentMethod {
		case "bank":
			ref := fmt.Sprintf("FSH-%d", order.ID)
			instructions := fiber.Map{
				"beneficiary_name":  "Furniture Shop Demo",
				"iban":              "BG00UNCR70001512345678",
				"bic":               "UNCRBGSF",
				"bank_name":         "UniCredit Bulbank (TEST)",
				"amount":            order.TotalPrice,
				"currency":          "EUR",
				"payment_reference": ref,
				"reason_line":       fmt.Sprintf("Order #%d - %s", order.ID, ref),
			}
			return c.JSON(fiber.Map{
				"order_id":                       order.ID,
				"total_price":                    order.TotalPrice,
				"estimated_production_time_days": order.EstimatedProductionTimeDays,
				"status":                         order.Status,
				"payment_status":                 order.PaymentStatus,
				"instructions":                   instructions,
			})
		case "card":
			fe := "http://localhost:5173"
			//Payment succeeds - 4242 4242 4242 4242
			//Payment requires authentication - 4000 0025 0000 3155
			//Payment is declined- 4000 0000 0000 9995
			stripe.Key = config.Env.StripeSecretKey

			amount := int64(math.Round(order.TotalPrice * 100))

			params := &stripe.CheckoutSessionParams{
				Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
				SuccessURL: stripe.String(fmt.Sprintf(
					"%s/payment/success?session_id={CHECKOUT_SESSION_ID}&order_id=%d",
					fe, order.ID,
				)),
				CancelURL: stripe.String(fmt.Sprintf(
					"%s/payment/cancel?order_id=%d",
					fe, order.ID,
				)),
				Metadata: map[string]string{"order_id": strconv.Itoa(int(order.ID))},
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
				return c.Status(500).JSON(fiber.Map{"message": "stripe error"})
			}

			return c.JSON(fiber.Map{
				"order_id":                       order.ID,
				"checkout_url":                   sess.URL,
				"estimated_production_time_days": order.EstimatedProductionTimeDays,
			})
		default:
			return c.JSON(fiber.Map{
				"order_id": order.ID, "total_price": order.TotalPrice,
				"estimated_production_time_days": order.EstimatedProductionTimeDays,
				"status":                         order.Status, "payment_status": order.PaymentStatus,
			})
		}
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

type patchStatusDTO struct {
	Status string `json:"status"`
}

func (h *Handler) AdminUpdateOrderStatus() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in patchStatusDTO
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
