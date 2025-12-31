package orders

import "github.com/gofiber/fiber/v2"

func RegisterCartRoutes(api fiber.Router, h *CartHandler) {
	grp := api.Group("/user")
	grp.Get("/cart", h.Get())
	grp.Put("/cart", h.Replace())
	grp.Post("/cart/items", h.AddItem())
	grp.Patch("/cart/items/:id", h.UpdateItem())
	grp.Delete("/cart/items/:id", h.RemoveItem())
	grp.Delete("/cart", h.Clear())
}
