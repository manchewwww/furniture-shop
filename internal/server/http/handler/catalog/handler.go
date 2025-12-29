package catalog

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	ec "furniture-shop/internal/entities/catalog"
	"furniture-shop/internal/service"
)

type Handler struct {
	svc service.CatalogService
}

func NewCatalogHandler(svc service.CatalogService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetDepartments() fiber.Handler {
	return func(c *fiber.Ctx) error {
		depts, err := h.svc.ListDepartments(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(depts)
	}
}

func (h *Handler) GetCategoriesByDepartment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
		cats, err := h.svc.ListCategoriesByDepartment(c.Context(), uint(id))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(cats)
	}
}

func (h *Handler) GetProductsByCategory() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
		products, err := h.svc.ListProductsByCategory(c.Context(), uint(id))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(products)
	}
}

func (h *Handler) GetProductDetails() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
		p, err := h.svc.GetProduct(c.Context(), uint(id))
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"message": "not found"})
		}
		return c.JSON(p)
	}
}

func (h *Handler) SearchProducts() fiber.Handler {
	return func(c *fiber.Ctx) error {
		q := c.Query("query")
		if q == "" {
			return c.JSON([]ec.Product{})
		}
		items, err := h.svc.SearchProducts(c.Context(), q, 50)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(items)
	}
}

func (h *Handler) GetProductRecommendations() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
		rec, err := h.svc.RecommendProducts(c.Context(), uint(id), 4)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(rec)
	}
}
