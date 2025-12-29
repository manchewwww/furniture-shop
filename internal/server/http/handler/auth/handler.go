package auth

import (
	"context"

	"github.com/gofiber/fiber/v2"

	eu "furniture-shop/internal/entities/user"
	"furniture-shop/internal/service"
	vld "furniture-shop/internal/validation"
)

type Handler struct {
	svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) *Handler {
	return &Handler{svc: svc}
}

type registerDTO struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Address  string `json:"address" validate:"omitempty,min=5"`
	Phone    string `json:"phone" validate:"omitempty,phone"`
}

func (h *Handler) Register() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in registerDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
		}
		if err := vld.ValidateStruct(in); err != nil {
			return err
		}
		user := eu.User{Role: "client", Name: in.Name, Email: in.Email, Address: in.Address, Phone: in.Phone}
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

func (h *Handler) createUser(ctx context.Context, u *eu.User) error {
	return h.svc.CreateUser(ctx, u)
}

type loginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *Handler) Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in loginDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
		}
		if err := vld.ValidateStruct(in); err != nil {
			return err
		}
		user, err := h.svc.Authenticate(c.Context(), in.Email, in.Password)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"message": "invalid email or password"})
		}
		token, _ := h.svc.GenerateJWT(user)
		return c.JSON(fiber.Map{"token": token, "user": fiber.Map{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role}})
	}
}

func (h *Handler) Me() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"id": c.Locals("user_id"), "email": c.Locals("user_email"), "role": c.Locals("user_role")})
	}
}
