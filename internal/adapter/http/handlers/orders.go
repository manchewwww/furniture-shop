package handlers

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
    app "furniture-shop/internal/app"
)

type OrdersHandler struct { svc app.OrdersService }

func NewOrdersHandler(svc app.OrdersService) *OrdersHandler { return &OrdersHandler{svc: svc} }

type orderItemIn struct {
    ProductID uint                     `json:"product_id"`
    Quantity  int                      `json:"quantity"`
    Options   []app.SelectedOption `json:"options"`
}
type createOrderDTO struct {
    UserID        *uint         `json:"user_id"`
    Name          string        `json:"name"`
    Email         string        `json:"email"`
    Address       string        `json:"address"`
    Phone         string        `json:"phone"`
    Items         []orderItemIn `json:"items"`
    PaymentMethod string        `json:"payment_method"`
}

func (h *OrdersHandler) CreateOrder() fiber.Handler {
    return func(c *fiber.Ctx) error {
        var in createOrderDTO
        if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"invalid request"}) }
        items := make([]app.CreateOrderItem, 0, len(in.Items))
        for _, it := range in.Items { items = append(items, app.CreateOrderItem{ProductID: it.ProductID, Quantity: it.Quantity, Options: it.Options}) }
        order, err := h.svc.CreateOrder(c.Context(), app.CreateOrderInput{
            UserID: in.UserID, Name: in.Name, Email: in.Email, Address: in.Address, Phone: in.Phone,
            Items: items, PaymentMethod: in.PaymentMethod,
        })
        if err != nil { return c.Status(400).JSON(fiber.Map{"message": err.Error()}) }
        return c.JSON(fiber.Map{
            "order_id": order.ID, "total_price": order.TotalPrice,
            "estimated_production_time_days": order.EstimatedProductionTimeDays,
            "status": order.Status, "payment_status": order.PaymentStatus,
        })
    }
}

func (h *OrdersHandler) UserOrders() fiber.Handler {
    return func(c *fiber.Ctx) error {
        uid, ok := c.Locals("user_id").(uint)
        if !ok { return c.Status(401).JSON(fiber.Map{"message":"unauthorized"}) }
        orders, err := h.svc.ListUserOrders(c.Context(), uid)
        if err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }
        return c.JSON(orders)
    }
}

func (h *OrdersHandler) UserOrderDetails() fiber.Handler {
    return func(c *fiber.Ctx) error {
        uid, ok := c.Locals("user_id").(uint)
        if !ok { return c.Status(401).JSON(fiber.Map{"message":"unauthorized"}) }
        // parse id
        var id uint
        _, _ = fmt.Sscan(c.Params("id"), &id)
        order, err := h.svc.GetUserOrder(c.Context(), uid, id)
        if err != nil { return c.Status(404).JSON(fiber.Map{"message":"not found"}) }
        return c.JSON(order)
    }
}

func (h *OrdersHandler) AdminListOrders() fiber.Handler {
    return func(c *fiber.Ctx) error {
        status := c.Query("status")
        orders, err := h.svc.AdminListOrders(c.Context(), status)
        if err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }
        return c.JSON(orders)
    }
}

type patchStatusDTO struct { Status string `json:"status"` }

func (h *OrdersHandler) AdminUpdateOrderStatus() fiber.Handler {
    return func(c *fiber.Ctx) error {
        var in patchStatusDTO
        if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"invalid request"}) }
        var id uint
        _, _ = fmt.Sscan(c.Params("id"), &id)
        if err := h.svc.AdminUpdateOrderStatus(c.Context(), id, in.Status); err != nil { return c.Status(400).JSON(fiber.Map{"message": err.Error()}) }
        return c.JSON(fiber.Map{"message":"updated"})
    }
}







