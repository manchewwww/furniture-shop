package domain

import (
    "furniture-shop/internal/service"
    "furniture-shop/internal/storage"
    sa "furniture-shop/internal/service/domain/auth"
    sadm "furniture-shop/internal/service/domain/admin"
    sc "furniture-shop/internal/service/domain/catalog"
    so "furniture-shop/internal/service/domain/orders"
    sp "furniture-shop/internal/service/domain/payments"
)

// NewService wires concrete domain services from repositories
func NewService(repos *storage.Repository, jwtSecret string) *service.Service {
    return &service.Service{
        Auth:    sa.NewAuthService(repos.Users, jwtSecret),
        Catalog: sc.NewCatalogService(repos.Departments, repos.Categories, repos.Products),
        Orders:  so.NewOrdersService(repos.Users, repos.Orders, repos.Products),
        Admin:   sadm.NewAdminService(repos.Departments, repos.Categories, repos.Products, repos.ProductOptions),
        Payment: sp.NewPaymentService(repos.Orders),
    }
}
