package handlers

import (
    "context"
    "github.com/gofiber/fiber/v2"

    models "furniture-shop/internal/domain/entity"
    "furniture-shop/internal/services"
)

type AuthHandler struct {
    svc services.AuthService
}

func NewAuthHandler(svc services.AuthService) *AuthHandler { return &AuthHandler{svc: svc} }

type registerDTO struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Address  string `json:"address"`
    Phone    string `json:"phone"`
}

func (h *AuthHandler) Register() fiber.Handler {
    return func(c *fiber.Ctx) error {
        var in registerDTO
        if err := c.BodyParser(&in); err != nil {
            return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
        }
        if in.Email == "" || in.Password == "" || in.Name == "" {
            return c.Status(400).JSON(fiber.Map{"message":"name, email and password are required"})
        }
        user := models.User{Role: "client", Name: in.Name, Email: in.Email, Address: in.Address, Phone: in.Phone}
        if err := user.SetPassword(in.Password); err != nil {
            return c.Status(500).JSON(fiber.Map{"message": "password hashing failed"})
        }
        if err := h.createUser(c.Context(), &user); err != nil {
            return c.Status(400).JSON(fiber.Map{"message": "could not create user"})
        }
        token, _ := h.svc.GenerateJWT(&user)
        return c.JSON(fiber.Map{"token": token, "user": fiber.Map{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role}})
    }
}

// createUser used to keep handler decoupled from infra; we rely on AuthService's repository
func (h *AuthHandler) createUser(ctx context.Context, u *models.User) error {
    return h.svc.CreateUser(ctx, u)
}

type loginDTO struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (h *AuthHandler) Login() fiber.Handler {
    return func(c *fiber.Ctx) error {
        var in loginDTO
        if err := c.BodyParser(&in); err != nil {
            return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
        }
        user, err := h.svc.Authenticate(c.Context(), in.Email, in.Password)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"message": "invalid email or password"})
        }
        token, _ := h.svc.GenerateJWT(user)
        return c.JSON(fiber.Map{"token": token, "user": fiber.Map{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role}})
    }
}

func (h *AuthHandler) Me() fiber.Handler {
    return func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "id":    c.Locals("user_id"),
            "email": c.Locals("user_email"),
            "role":  c.Locals("user_role"),
        })
    }
}
