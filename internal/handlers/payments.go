package handlers

import (
    "regexp"
    "github.com/gofiber/fiber/v2"

    "furniture-shop/internal/database"
    "furniture-shop/internal/models"
)

type cardDTO struct {
    OrderID       uint   `json:"order_id"`
    Cardholder    string `json:"cardholder_name"`
    CardNumber    string `json:"card_number"`
    ExpiryMonth   string `json:"expiry_month"`
    ExpiryYear    string `json:"expiry_year"`
    CVV           string `json:"cvv"`
}

func PayByCard(c *fiber.Ctx) error {
    var in cardDTO
    if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"Невалидни данни"}) }
    var order models.Order
    if err := database.DB.First(&order, in.OrderID).Error; err != nil { return c.Status(404).JSON(fiber.Map{"message":"Поръчката не е намерена"}) }
    if order.PaymentMethod != "карта" { return c.Status(400).JSON(fiber.Map{"message":"Методът на плащане не е карта"}) }

    if simulateCreditCardAuthorization(in) {
        database.DB.Model(&order).Update("payment_status", "платено")
        return c.JSON(fiber.Map{"message":"Плащането е успешно", "payment_status":"платено"})
    }
    database.DB.Model(&order).Update("payment_status", "отказано")
    return c.Status(402).JSON(fiber.Map{"message":"Плащането е отказано", "payment_status":"отказано"})
}

func simulateCreditCardAuthorization(in cardDTO) bool {
    reDigits := regexp.MustCompile(`^\d+$`)
    if !reDigits.MatchString(in.CardNumber) || len(in.CardNumber) < 12 || len(in.CardNumber) > 19 { return false }
    if !reDigits.MatchString(in.CVV) || (len(in.CVV) != 3 && len(in.CVV) != 4) { return false }
    if in.Cardholder == "" { return false }
    if !reDigits.MatchString(in.ExpiryMonth) || !reDigits.MatchString(in.ExpiryYear) { return false }
    return true
}

