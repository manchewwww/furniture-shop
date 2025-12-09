package handlers

import (
    "github.com/gofiber/fiber/v2"

    "furniture-shop/internal/config"
    "furniture-shop/internal/database"
    "furniture-shop/internal/models"
    "furniture-shop/internal/services"
)

type registerDTO struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Address  string `json:"address"`
    Phone    string `json:"phone"`
}

func Register(cfg *config.Config) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var in registerDTO
        if err := c.BodyParser(&in); err != nil {
            return c.Status(400).JSON(fiber.Map{"message": "Невалидни данни"})
        }
        if in.Email == "" || in.Password == "" || in.Name == "" {
            return c.Status(400).JSON(fiber.Map{"message":"Име, имейл и парола са задължителни"})
        }
        user := models.User{Role: "client", Name: in.Name, Email: in.Email, Address: in.Address, Phone: in.Phone}
        if err := user.SetPassword(in.Password); err != nil {
            return c.Status(500).JSON(fiber.Map{"message": "Грешка при регистрация"})
        }
        if err := database.DB.Create(&user).Error; err != nil {
            return c.Status(400).JSON(fiber.Map{"message": "Имейлът вече съществува"})
        }
        token, _ := services.GenerateJWT(&user, cfg)
        return c.JSON(fiber.Map{"token": token, "user": fiber.Map{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role}})
    }
}

type loginDTO struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func Login(cfg *config.Config) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var in loginDTO
        if err := c.BodyParser(&in); err != nil {
            return c.Status(400).JSON(fiber.Map{"message": "Невалидни данни"})
        }
        user, err := services.Authenticate(in.Email, in.Password)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"message": "Невалиден имейл или парола"})
        }
        token, _ := services.GenerateJWT(user, cfg)
        return c.JSON(fiber.Map{"token": token, "user": fiber.Map{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role}})
    }
}

func Me() fiber.Handler {
    return func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "id":    c.Locals("user_id"),
            "email": c.Locals("user_email"),
            "role":  c.Locals("user_role"),
        })
    }
}

