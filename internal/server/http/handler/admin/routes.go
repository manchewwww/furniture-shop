package admin

import (
	ho "furniture-shop/internal/server/http/handler/orders"

	"github.com/gofiber/fiber/v2"
)

// Register admin-specific routes; orders admin endpoints reuse orders handler
func RegisterAdminRoutes(admin fiber.Router, h *Handler, orders *ho.Handler) {
	admin.Get("/departments", h.ListDepartments())
	admin.Post("/departments", h.CreateDepartment())
	admin.Put("/departments/:id", h.UpdateDepartment())
	admin.Delete("/departments/:id", h.DeleteDepartment())

	admin.Get("/categories", h.ListCategories())
	admin.Post("/categories", h.CreateCategory())
	admin.Put("/categories/:id", h.UpdateCategory())
	admin.Delete("/categories/:id", h.DeleteCategory())

	admin.Get("/products", h.ListProducts())
	admin.Post("/products", h.CreateProduct())
	admin.Put("/products/:id", h.UpdateProduct())
	admin.Delete("/products/:id", h.DeleteProduct())

	admin.Get("/product_options", h.ListProductOptions())
	admin.Post("/product_options", h.CreateProductOption())
	admin.Put("/product_options/:id", h.UpdateProductOption())
	admin.Delete("/product_options/:id", h.DeleteProductOption())
	admin.Post("/upload", h.UploadImage())

	admin.Get("/orders", orders.AdminListOrders())
	admin.Patch("/orders/:id/status", orders.AdminUpdateOrderStatus())
}
