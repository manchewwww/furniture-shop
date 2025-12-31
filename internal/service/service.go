package service

import (
	"context"

	"furniture-shop/internal/dtos/cart"
	order_dto "furniture-shop/internal/dtos/orders"
	ec "furniture-shop/internal/entities/catalog"
	ei "furniture-shop/internal/entities/inventory"
	eo "furniture-shop/internal/entities/orders"
	eu "furniture-shop/internal/entities/user"
)

type AuthService interface {
	GenerateJWT(u *eu.User) (string, error)
	Authenticate(ctx context.Context, email, password string) (*eu.User, error)
	CreateUser(ctx context.Context, u *eu.User) error
}

type CatalogService interface {
	ListDepartments(ctx context.Context) ([]ec.Department, error)
	ListCategoriesByDepartment(ctx context.Context, departmentID uint) ([]ec.Category, error)
	ListProductsByCategory(ctx context.Context, categoryID uint) ([]ec.Product, error)
	GetProduct(ctx context.Context, id uint) (*ec.Product, error)
	SearchProducts(ctx context.Context, query string, limit int) ([]ec.Product, error)
	RecommendProducts(ctx context.Context, productID uint, limit int) ([]ec.Product, error)
}

type OrdersService interface {
	CreateOrder(ctx context.Context, in order_dto.CreateOrderInput) (*eo.Order, error)
	ListUserOrders(ctx context.Context, userID uint) ([]eo.Order, error)
	GetUserOrder(ctx context.Context, userID, orderID uint) (*eo.Order, error)
	AdminListOrders(ctx context.Context, status string) ([]eo.Order, error)
	AdminUpdateOrderStatus(ctx context.Context, orderID uint, status string) error
}

type AdminService interface {
	ListDepartments(ctx context.Context) ([]ec.Department, error)
	CreateDepartment(ctx context.Context, d *ec.Department) error
	UpdateDepartment(ctx context.Context, id uint, d ec.Department) error
	DeleteDepartment(ctx context.Context, id uint) error
	ListCategories(ctx context.Context) ([]ec.Category, error)
	CreateCategory(ctx context.Context, c *ec.Category) error
	UpdateCategory(ctx context.Context, id uint, c ec.Category) error
	DeleteCategory(ctx context.Context, id uint) error
	ListProducts(ctx context.Context) ([]ec.Product, error)
	CreateProduct(ctx context.Context, p *ec.Product) error
	UpdateProduct(ctx context.Context, id uint, p ec.Product) error
	DeleteProduct(ctx context.Context, id uint) error
	ListProductOptions(ctx context.Context, productID *uint) ([]ec.ProductOption, error)
	CreateProductOption(ctx context.Context, o *ec.ProductOption) error
	UpdateProductOption(ctx context.Context, id uint, o ec.ProductOption) error
	DeleteProductOption(ctx context.Context, id uint) error
	ListStock(ctx context.Context) ([]ei.Stock, error)
	UpsertStock(ctx context.Context, material string, qty float64, unit string) error
}

type PaymentService interface {
	ProcessPaymentResult(ctx context.Context, orderID uint, paymentStatus, orderStatus string) error
}

type CartService interface {
	Get(ctx context.Context, userID uint) (*eo.Cart, error)
	Replace(ctx context.Context, userID uint, in cart.ReplaceCartRequest) (*eo.Cart, error)
	AddItem(ctx context.Context, userID uint, in cart.AddCartItemRequest) (*eo.CartItem, error)
	UpdateItem(ctx context.Context, userID uint, itemID uint, in cart.UpdateCartItemRequest) error
	RemoveItem(ctx context.Context, userID uint, itemID uint) error
	Clear(ctx context.Context, userID uint) error
}

type Service struct {
	Auth    AuthService
	Catalog CatalogService
	Orders  OrdersService
	Admin   AdminService
	Payment PaymentService
	Cart    CartService
}
