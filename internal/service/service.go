package service

import (
	"context"

	ec "furniture-shop/internal/entities/catalog"
	eo "furniture-shop/internal/entities/orders"
	eu "furniture-shop/internal/entities/user"
)

// Shared DTOs
type SelectedOption struct {
	ID   uint   `json:"id"`
	Type string `json:"type"`
}

type CreateOrderItem struct {
	ProductID uint             `json:"product_id"`
	Quantity  int              `json:"quantity"`
	Options   []SelectedOption `json:"options"`
}

type CreateOrderInput struct {
	UserID        *uint             `json:"user_id"`
	Name          string            `json:"name"`
	Email         string            `json:"email"`
	Address       string            `json:"address"`
	Phone         string            `json:"phone"`
	Items         []CreateOrderItem `json:"items"`
	PaymentMethod string            `json:"payment_method"`
}

type CardPayment struct {
	OrderID     uint   `json:"order_id" validate:"required,gt=0"`
	Cardholder  string `json:"cardholder_name" validate:"required,min=2"`
	CardNumber  string `json:"card_number" validate:"required,card"`
	ExpiryMonth string `json:"expiry_month" validate:"required,month"`
	ExpiryYear  string `json:"expiry_year" validate:"required,numeric"`
	CVV         string `json:"cvv" validate:"required,cvv"`
}

// Auth
type AuthService interface {
	GenerateJWT(u *eu.User) (string, error)
	Authenticate(ctx context.Context, email, password string) (*eu.User, error)
	CreateUser(ctx context.Context, u *eu.User) error
}

// Catalog
type CatalogService interface {
	ListDepartments(ctx context.Context) ([]ec.Department, error)
	ListCategoriesByDepartment(ctx context.Context, departmentID uint) ([]ec.Category, error)
	ListProductsByCategory(ctx context.Context, categoryID uint) ([]ec.Product, error)
	GetProduct(ctx context.Context, id uint) (*ec.Product, error)
	SearchProducts(ctx context.Context, query string, limit int) ([]ec.Product, error)
	RecommendProducts(ctx context.Context, productID uint, limit int) ([]ec.Product, error)
}

// Orders
type OrdersService interface {
	CreateOrder(ctx context.Context, in CreateOrderInput) (*eo.Order, error)
	ListUserOrders(ctx context.Context, userID uint) ([]eo.Order, error)
	GetUserOrder(ctx context.Context, userID, orderID uint) (*eo.Order, error)
	AdminListOrders(ctx context.Context, status string) ([]eo.Order, error)
	AdminUpdateOrderStatus(ctx context.Context, orderID uint, status string) error
}

// Admin
type AdminService interface {
	// Departments
	ListDepartments(ctx context.Context) ([]ec.Department, error)
	CreateDepartment(ctx context.Context, d *ec.Department) error
	UpdateDepartment(ctx context.Context, id uint, d ec.Department) error
	DeleteDepartment(ctx context.Context, id uint) error
	// Categories
	ListCategories(ctx context.Context) ([]ec.Category, error)
	CreateCategory(ctx context.Context, c *ec.Category) error
	UpdateCategory(ctx context.Context, id uint, c ec.Category) error
	DeleteCategory(ctx context.Context, id uint) error
	// Products
	ListProducts(ctx context.Context) ([]ec.Product, error)
	CreateProduct(ctx context.Context, p *ec.Product) error
	UpdateProduct(ctx context.Context, id uint, p ec.Product) error
	DeleteProduct(ctx context.Context, id uint) error
	// Product Options
	ListProductOptions(ctx context.Context, productID *uint) ([]ec.ProductOption, error)
	CreateProductOption(ctx context.Context, o *ec.ProductOption) error
	UpdateProductOption(ctx context.Context, id uint, o ec.ProductOption) error
	DeleteProductOption(ctx context.Context, id uint) error
}

// Payments
type PaymentService interface {
	PayByCard(ctx context.Context, in CardPayment) (string, error)
}

// Service is the top-level container
type Service struct {
	Auth    AuthService
	Catalog CatalogService
	Orders  OrdersService
	Admin   AdminService
	Payment PaymentService
}
