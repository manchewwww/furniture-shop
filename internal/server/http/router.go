package http

import (
	"github.com/gofiber/fiber/v2"

	ha "furniture-shop/internal/server/http/handler/admin"
	hau "furniture-shop/internal/server/http/handler/auth"
	hc "furniture-shop/internal/server/http/handler/catalog"
	ho "furniture-shop/internal/server/http/handler/orders"
	hp "furniture-shop/internal/server/http/handler/payments"
	"furniture-shop/internal/server/http/middleware"
)

func buildRoutes(s *Server) {
	s.app.Get("/health", func(c *fiber.Ctx) error { return c.SendString("ok") })
	api := s.app.Group("/api")

	authH := hau.NewAuthHandler(s.svc.Auth)
	catalogH := hc.NewCatalogHandler(s.svc.Catalog)
	ordersH := ho.NewOrdersHandler(s.svc.Orders)
	adminH := ha.NewAdminHandler(s.svc.Admin)
	paymentsH := hp.NewPaymentsHandler(s.svc.Payment)

	// Auth
	hau.Register(api, authH)

	// Catalog
	hc.Register(api, catalogH)

	// Orders and payments
	api.Post("/orders", middleware.JWTAuth(), ordersH.CreateOrder())
	hp.Register(api, paymentsH)

	// Authenticated user routes
	authGroup := api.Group("/user", middleware.JWTAuth())
	authGroup.Get("/me", authH.Me())
	authGroup.Get("/orders", ordersH.UserOrders())
	authGroup.Get("/orders/:id", ordersH.UserOrderDetails())

	// Admin routes
	adminGroup := api.Group("/admin", middleware.JWTAuth(), middleware.RequireAdmin)
	ha.RegisterAdminRoutes(adminGroup, adminH, ordersH)
}
