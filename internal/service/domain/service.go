package domain

import (
	"furniture-shop/internal/service"
	sadm "furniture-shop/internal/service/domain/admin"
	sa "furniture-shop/internal/service/domain/auth"
	sc "furniture-shop/internal/service/domain/catalog"
	so "furniture-shop/internal/service/domain/orders"
	sp "furniture-shop/internal/service/domain/payments"
	"furniture-shop/internal/service/mailer"
	"furniture-shop/internal/storage"
)

// NewService wires concrete domain services from repositories
func NewService(repos *storage.Repository, jwtSecret string) *service.Service {
	return &service.Service{
		Auth:    sa.NewAuthService(repos.Users, jwtSecret),
		Catalog: sc.NewCatalogService(repos.Departments, repos.Categories, repos.Products),
		Orders:  so.NewOrdersService(repos.Users, repos.Orders, repos.Products, repos.Stock),
		Admin:   sadm.NewAdminService(repos.Departments, repos.Categories, repos.Products, repos.ProductOptions, repos.Stock),
		Payment: sp.NewPaymentService(repos.Orders, repos.Products, repos.Stock, repos.Users, mailer.NewSender()),
		Cart:    so.NewCartService(repos.Carts),
	}
}
