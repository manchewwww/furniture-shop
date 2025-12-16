package handlers

import (
    "github.com/gofiber/fiber/v2"
    "furniture-shop/internal/services"
)

type PaymentsHandler struct { svc services.PaymentService }

func NewPaymentsHandler(svc services.PaymentService) *PaymentsHandler { return &PaymentsHandler{svc: svc} }

type cardDTO struct {
    OrderID     uint   `json:"order_id"`
    Cardholder  string `json:"cardholder_name"`
    CardNumber  string `json:"card_number"`
    ExpiryMonth string `json:"expiry_month"`
    ExpiryYear  string `json:"expiry_year"`
    CVV         string `json:"cvv"`
}

func (h *PaymentsHandler) PayByCard() fiber.Handler {
    return func(c *fiber.Ctx) error {
        var in cardDTO
        if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"invalid request"}) }
        status, err := h.svc.PayByCard(c.Context(), services.CardPayment{
            OrderID: in.OrderID, Cardholder: in.Cardholder, CardNumber: in.CardNumber, ExpiryMonth: in.ExpiryMonth, ExpiryYear: in.ExpiryYear, CVV: in.CVV,
        })
        if err != nil { return c.Status(402).JSON(fiber.Map{"message":"payment failed", "payment_status":"declined"}) }
        return c.JSON(fiber.Map{"message":"payment accepted", "payment_status": status})
    }
}
