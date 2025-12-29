package admin

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"

	ec "furniture-shop/internal/entities/catalog"
	"furniture-shop/internal/service"
	vld "furniture-shop/internal/validation"
)

type DepartmentDTO struct {
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description" validate:"omitempty,min=2"`
	ImageURL    string `json:"image_url" validate:"omitempty,url"`
}

type CategoryDTO struct {
	DepartmentID uint   `json:"department_id" validate:"required,gt=0"`
	Name         string `json:"name" validate:"required,min=2"`
	Description  string `json:"description" validate:"omitempty,min=2"`
}

type ProductDTO struct {
	CategoryID             uint    `json:"category_id" validate:"required,gt=0"`
	Name                   string  `json:"name" validate:"required,min=2"`
	ShortDescription       string  `json:"short_description" validate:"omitempty,min=2"`
	LongDescription        string  `json:"long_description" validate:"omitempty,min=2"`
	BasePrice              float64 `json:"base_price" validate:"required,gte=0"`
	BaseProductionTimeDays int     `json:"base_production_time_days" validate:"required,gte=0"`
	ImageURL               string  `json:"image_url" validate:"omitempty,url"`
}

type ProductOptionDTO struct {
	ProductID                     uint    `json:"product_id" validate:"required,gt=0"`
	OptionType                    string  `json:"option_type" validate:"required,oneof=color size material extra"`
	OptionName                    string  `json:"option_name" validate:"required,min=1"`
	PriceModifierType             string  `json:"price_modifier_type" validate:"required,oneof=absolute percent"`
	PriceModifierValue            float64 `json:"price_modifier_value" validate:"required"`
	ProductionTimeModifierDays    int     `json:"production_time_modifier_days" validate:"omitempty"`
	ProductionTimeModifierPercent *int    `json:"production_time_modifier_percent" validate:"omitempty"`
}

type Handler struct {
	svc service.AdminService
}

func NewAdminHandler(svc service.AdminService) *Handler {
	return &Handler{svc: svc}
}

// Departments
func (h *Handler) ListDepartments() fiber.Handler {
	return func(c *fiber.Ctx) error {
		items, err := h.svc.ListDepartments(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(items)
	}
}

func (h *Handler) CreateDepartment() fiber.Handler {
	return func(c *fiber.Ctx) error {

		var in DepartmentDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
		}
		if err := vld.ValidateStruct(in); err != nil {
			return err
		}
		dep := ec.Department{Name: in.Name, Description: in.Description, ImageURL: in.ImageURL}
		if err := h.svc.CreateDepartment(c.Context(), &dep); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(dep)
	}
}

func (h *Handler) UpdateDepartment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in DepartmentDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
		}
		if err := vld.ValidateStruct(in); err != nil {
			return err
		}
		var id uint
		_, _ = fmt.Sscan(c.Params("id"), &id)
		if err := h.svc.UpdateDepartment(c.Context(), id, ec.Department{Name: in.Name, Description: in.Description, ImageURL: in.ImageURL}); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(fiber.Map{"message": "updated"})
	}
}

func (h *Handler) DeleteDepartment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var id uint
		if _, err := fmt.Sscan(c.Params("id"), &id); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid id"})
		}
		if err := h.svc.DeleteDepartment(c.Context(), id); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(fiber.Map{"message": "deleted"})
	}
}

// Categories
func (h *Handler) ListCategories() fiber.Handler {
	return func(c *fiber.Ctx) error {
		items, err := h.svc.ListCategories(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(items)
	}
}

func (h *Handler) CreateCategory() fiber.Handler {
	return func(c *fiber.Ctx) error {

		var in CategoryDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
		}
		if err := vld.ValidateStruct(in); err != nil {
			return err
		}
		cat := ec.Category{DepartmentID: in.DepartmentID, Name: in.Name, Description: in.Description}
		if err := h.svc.CreateCategory(c.Context(), &cat); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(cat)
	}
}

func (h *Handler) UpdateCategory() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in CategoryDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
		}
		if err := vld.ValidateStruct(in); err != nil {
			return err
		}
		var id uint
		if _, err := fmt.Sscan(c.Params("id"), &id); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid id"})
		}
		if err := h.svc.UpdateCategory(c.Context(), id, ec.Category{DepartmentID: in.DepartmentID, Name: in.Name, Description: in.Description}); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(fiber.Map{"message": "updated"})
	}
}

func (h *Handler) DeleteCategory() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var id uint
		if _, err := fmt.Sscan(c.Params("id"), &id); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid id"})
		}
		if err := h.svc.DeleteCategory(c.Context(), id); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(fiber.Map{"message": "deleted"})
	}
}

// Products
func (h *Handler) ListProducts() fiber.Handler {
	return func(c *fiber.Ctx) error {
		items, err := h.svc.ListProducts(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(items)
	}
}

func (h *Handler) CreateProduct() fiber.Handler {
	return func(c *fiber.Ctx) error {

		var in ProductDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
		}
		if err := vld.ValidateStruct(in); err != nil {
			return err
		}
		p := ec.Product{
			CategoryID:             in.CategoryID,
			Name:                   in.Name,
			ShortDescription:       in.ShortDescription,
			LongDescription:        in.LongDescription,
			BasePrice:              in.BasePrice,
			BaseProductionTimeDays: in.BaseProductionTimeDays,
			ImageURL:               in.ImageURL,
		}
		if err := h.svc.CreateProduct(c.Context(), &p); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(p)
	}
}

func (h *Handler) UpdateProduct() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in ProductDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
		}
		if err := vld.ValidateStruct(in); err != nil {
			return err
		}
		var id uint
		if _, err := fmt.Sscan(c.Params("id"), &id); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid id"})
		}
		if err := h.svc.UpdateProduct(c.Context(), id, ec.Product{
			CategoryID:             in.CategoryID,
			Name:                   in.Name,
			ShortDescription:       in.ShortDescription,
			LongDescription:        in.LongDescription,
			BasePrice:              in.BasePrice,
			BaseProductionTimeDays: in.BaseProductionTimeDays,
			ImageURL:               in.ImageURL,
		}); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(fiber.Map{"message": "updated"})
	}
}

func (h *Handler) DeleteProduct() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var id uint
		if _, err := fmt.Sscan(c.Params("id"), &id); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid id"})
		}
		if err := h.svc.DeleteProduct(c.Context(), id); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(fiber.Map{"message": "deleted"})
	}
}

// Product Options
func (h *Handler) ListProductOptions() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var pid *uint
		if s := c.Query("product_id"); s != "" {
			var v uint
			if _, err := fmt.Sscan(s, &v); err != nil {
				return c.Status(400).JSON(fiber.Map{"message": "invalid product_id"})
			}
			pid = &v
		}
		items, err := h.svc.ListProductOptions(c.Context(), pid)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(items)
	}
}

func (h *Handler) CreateProductOption() fiber.Handler {
	return func(c *fiber.Ctx) error {

		var in ProductOptionDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
		}
		if err := vld.ValidateStruct(in); err != nil {
			return err
		}
		po := ec.ProductOption{
			ProductID:                     in.ProductID,
			OptionType:                    in.OptionType,
			OptionName:                    in.OptionName,
			PriceModifierType:             in.PriceModifierType,
			PriceModifierValue:            in.PriceModifierValue,
			ProductionTimeModifierDays:    in.ProductionTimeModifierDays,
			ProductionTimeModifierPercent: in.ProductionTimeModifierPercent,
		}
		if err := h.svc.CreateProductOption(c.Context(), &po); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(po)
	}
}

func (h *Handler) UpdateProductOption() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in ProductOptionDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
		}
		if err := vld.ValidateStruct(in); err != nil {
			return err
		}
		var id uint
		if _, err := fmt.Sscan(c.Params("id"), &id); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid id"})
		}
		if err := h.svc.UpdateProductOption(c.Context(), id, ec.ProductOption{
			ProductID:                     in.ProductID,
			OptionType:                    in.OptionType,
			OptionName:                    in.OptionName,
			PriceModifierType:             in.PriceModifierType,
			PriceModifierValue:            in.PriceModifierValue,
			ProductionTimeModifierDays:    in.ProductionTimeModifierDays,
			ProductionTimeModifierPercent: in.ProductionTimeModifierPercent,
		}); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(fiber.Map{"message": "updated"})
	}
}

func (h *Handler) DeleteProductOption() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var id uint
		if _, err := fmt.Sscan(c.Params("id"), &id); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid id"})
		}
		if err := h.svc.DeleteProductOption(c.Context(), id); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		return c.JSON(fiber.Map{"message": "deleted"})
	}
}

// UploadImage handles admin image uploads and returns a public URL
func (h *Handler) UploadImage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "file is required"})
		}
		if err := os.MkdirAll("uploads", 0o755); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "server error"})
		}
		ext := filepath.Ext(fileHeader.Filename)
		if ext == "" {
			ext = ".bin"
		}
		name := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), "upload", ext)
		dst := filepath.Join("uploads", name)
		if err := c.SaveFile(fileHeader, dst); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "save failed"})
		}
		return c.JSON(fiber.Map{"url": "/uploads/" + name})
	}
}
