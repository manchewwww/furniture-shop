package handlers

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"

    "furniture-shop/internal/database"
    "furniture-shop/internal/models"
    "furniture-shop/internal/services"
)

type orderItemIn struct {
    ProductID uint                    `json:"product_id"`
    Quantity  int                     `json:"quantity"`
    Options   []services.SelectedOption `json:"options"`
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

func CreateOrder(c *fiber.Ctx) error {
    var in createOrderDTO
    if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"Невалидни данни"}) }
    if len(in.Items) == 0 { return c.Status(400).JSON(fiber.Map{"message":"Празна количка"}) }
    if in.PaymentMethod == "" { return c.Status(400).JSON(fiber.Map{"message":"Изберете метод на плащане"}) }

    // Ensure user
    var user models.User
    if in.UserID != nil && *in.UserID != 0 {
        if err := database.DB.First(&user, *in.UserID).Error; err != nil { return c.Status(400).JSON(fiber.Map{"message":"Невалиден потребител"}) }
    } else {
        // create guest user record (optional). Here we reuse/create by email if provided.
        if in.Email == "" { return c.Status(400).JSON(fiber.Map{"message":"Имейл е задължителен за поръчка"}) }
        if err := database.DB.Where("email = ?", in.Email).First(&user).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                user = models.User{Role: "client", Name: in.Name, Email: in.Email, Address: in.Address, Phone: in.Phone}
                _ = user.SetPassword("guest")
                if err := database.DB.Create(&user).Error; err != nil { return c.Status(500).JSON(fiber.Map{"message":"Грешка при създаване на потребител"}) }
            } else { return c.Status(500).JSON(fiber.Map{"message":"Грешка"}) }
        }
    }

    order := models.Order{UserID: user.ID, Status: "нова", PaymentMethod: in.PaymentMethod, PaymentStatus: "чакаплащане"}
    var items []models.OrderItem
    var total float64
    for _, it := range in.Items {
        var p models.Product
        if err := database.DB.Preload("Options").First(&p, it.ProductID).Error; err != nil {
            return c.Status(400).JSON(fiber.Map{"message": fmt.Sprintf("Продукт %d не е намерен", it.ProductID)})
        }
        if it.Quantity <= 0 { it.Quantity = 1 }
        unit := services.CalculateUnitPrice(p, it.Options)
        line := unit * float64(it.Quantity)
        pt := services.CalculateItemProductionTime(p, it.Options)
        items = append(items, models.OrderItem{
            ProductID: p.ID,
            Quantity:  it.Quantity,
            UnitPrice: unit,
            LineTotal: line,
            CalculatedProductionTimeDays: pt,
            SelectedOptionsJSON: services.MarshalSelectedOptions(it.Options),
        })
        total += line
    }
    order.TotalPrice = total
    // preliminary to compute overall production
    order.Items = items
    order.EstimatedProductionTimeDays = services.CalculateOrderProductionTime(items)

    if err := database.DB.Create(&order).Error; err != nil { return c.Status(500).JSON(fiber.Map{"message":"Грешка при запис"}) }
    for i := range items { items[i].OrderID = order.ID }
    if err := database.DB.Create(&items).Error; err != nil { return c.Status(500).JSON(fiber.Map{"message":"Грешка при запис на позиции"}) }

    // Simple virtual email log
    // log.Printf("Нова поръчка #%d за %s", order.ID, user.Email)

    return c.JSON(fiber.Map{
        "order_id": order.ID, "total_price": order.TotalPrice,
        "estimated_production_time_days": order.EstimatedProductionTimeDays,
        "status": order.Status, "payment_status": order.PaymentStatus,
    })
}

func UserOrders(c *fiber.Ctx) error {
    uid, ok := c.Locals("user_id").(uint)
    if !ok { return c.Status(401).JSON(fiber.Map{"message":"Неоторизиран"}) }
    var orders []models.Order
    if err := database.DB.Where("user_id = ?", uid).Order("created_at DESC").Find(&orders).Error; err != nil { return c.Status(500).JSON(fiber.Map{"message":"Грешка"}) }
    return c.JSON(orders)
}

func UserOrderDetails(c *fiber.Ctx) error {
    uid, ok := c.Locals("user_id").(uint)
    if !ok { return c.Status(401).JSON(fiber.Map{"message":"Неоторизиран"}) }
    var order models.Order
    if err := database.DB.Preload("Items").First(&order, c.Params("id")).Error; err != nil { return c.Status(404).JSON(fiber.Map{"message":"Не е намерена"}) }
    if order.UserID != uid { return c.Status(403).JSON(fiber.Map{"message":"Забранено"}) }
    return c.JSON(order)
}

func AdminListOrders(c *fiber.Ctx) error {
    status := c.Query("status")
    var orders []models.Order
    q := database.DB.Order("created_at DESC")
    if status != "" { q = q.Where("status = ?", status) }
    if err := q.Find(&orders).Error; err != nil { return c.Status(500).JSON(fiber.Map{"message":"Грешка"}) }
    return c.JSON(orders)
}

type patchStatusDTO struct { Status string `json:"status"` }

func AdminUpdateOrderStatus(c *fiber.Ctx) error {
    var in patchStatusDTO
    if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"Невалидни данни"}) }
    allowed := map[string]bool{"нова":true,"потвърдена":true,"впроизводство":true,"изпратена":true,"доставена":true,"отказана":true}
    if !allowed[in.Status] { return c.Status(400).JSON(fiber.Map{"message":"Невалиден статус"}) }
    if err := database.DB.Model(&models.Order{}).Where("id = ?", c.Params("id")).Update("status", in.Status).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"message":"Грешка"})
    }
    return c.JSON(fiber.Map{"message":"Статусът е обновен"})
}

