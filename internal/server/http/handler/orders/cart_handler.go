package orders

import (
	"fmt"

	cartdto "furniture-shop/internal/dtos/cart"
	"furniture-shop/internal/service"

	"github.com/gofiber/fiber/v2"
)

type CartHandler struct {
	svc service.CartService
}

func NewCartHandler(svc service.CartService) *CartHandler {
	return &CartHandler{svc: svc}
}

func userID(c *fiber.Ctx) (uint, bool) {
	uid, ok := c.Locals("user_id").(uint)
	return uid, ok
}

func (h *CartHandler) Get() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := userID(c)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"message": "unauthorized"})
		}
		cart, err := h.svc.Get(c.Context(), uid)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(cart)
	}
}

func (h *CartHandler) Replace() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := userID(c)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"message": "unauthorized"})
		}
		var in cartdto.ReplaceCartRequest
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
		}
		cart, err := h.svc.Replace(c.Context(), uid, in)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(cart)
	}
}

func (h *CartHandler) AddItem() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := userID(c)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"message": "unauthorized"})
		}
		var in cartdto.AddCartItemRequest
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
		}
		item, err := h.svc.AddItem(c.Context(), uid, in)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": fmt.Sprintf("%v", err)})
		}
		return c.JSON(item)
	}
}

func (h *CartHandler) UpdateItem() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := userID(c)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"message": "unauthorized"})
		}
		var id uint
		if _, err := fmt.Sscan(c.Params("id"), &id); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid id"})
		}
		var in cartdto.UpdateCartItemRequest
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
		}
		if err := h.svc.UpdateItem(c.Context(), uid, id, in); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(fiber.Map{"message": "updated"})
	}
}

func (h *CartHandler) RemoveItem() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := userID(c)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"message": "unauthorized"})
		}
		var id uint
		if _, err := fmt.Sscan(c.Params("id"), &id); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid id"})
		}
		if err := h.svc.RemoveItem(c.Context(), uid, id); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(fiber.Map{"message": "deleted"})
	}
}

func (h *CartHandler) Clear() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := userID(c)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"message": "unauthorized"})
		}
		if err := h.svc.Clear(c.Context(), uid); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(fiber.Map{"message": "cleared"})
	}
}
