package orders

import "github.com/gofiber/fiber/v2"

func RegisterCartRoutes(r fiber.Router, h *CartHandler) {
	r.Get("/cart", h.Get())
	r.Put("/cart", h.Replace())
	r.Post("/cart/items", h.AddItem())
	r.Patch("/cart/items/:id", h.UpdateItem())
	r.Delete("/cart/items/:id", h.RemoveItem())
	r.Delete("/cart", h.Clear())
}
