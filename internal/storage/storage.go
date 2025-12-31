package storage

import (
	"context"

	ec "furniture-shop/internal/entities/catalog"
	eo "furniture-shop/internal/entities/orders"
	eu "furniture-shop/internal/entities/user"
)

// User
type UserRepository interface {
	Create(ctx context.Context, u *eu.User) error
	FindByEmail(ctx context.Context, email string) (*eu.User, error)
	FindByID(ctx context.Context, id uint) (*eu.User, error)
}

// Catalog
type DepartmentRepository interface {
	List(ctx context.Context) ([]ec.Department, error)
	Create(ctx context.Context, d *ec.Department) error
	Update(ctx context.Context, id uint, d ec.Department) error
	Delete(ctx context.Context, id uint) error
}

// Category
type CategoryRepository interface {
	ListByDepartment(ctx context.Context, departmentID uint) ([]ec.Category, error)
	ListAll(ctx context.Context) ([]ec.Category, error)
	Create(ctx context.Context, c *ec.Category) error
	Update(ctx context.Context, id uint, c ec.Category) error
	Delete(ctx context.Context, id uint) error
}

type ProductRepository interface {
	ListByCategory(ctx context.Context, categoryID uint) ([]ec.Product, error)
	FindByID(ctx context.Context, id uint) (*ec.Product, error)
	Search(ctx context.Context, query string, limit int) ([]ec.Product, error)
	ListRecommendations(ctx context.Context, p *ec.Product, limit int) ([]ec.Product, error)
	IncrementRecommendation(ctx context.Context, productID uint) error
	ListAll(ctx context.Context) ([]ec.Product, error)
	Create(ctx context.Context, p *ec.Product) error
	Update(ctx context.Context, id uint, p ec.Product) error
	Delete(ctx context.Context, id uint) error
	AdjustQuantity(ctx context.Context, productID uint, delta int) error
}

// Cart persistence
type CartRepository interface {
	GetOrCreateByUser(ctx context.Context, userID uint) (*eo.Cart, error)
	ReplaceItems(ctx context.Context, userID uint, items []eo.CartItem) (*eo.Cart, error)
	AddItem(ctx context.Context, userID uint, item *eo.CartItem) (*eo.CartItem, error)
	UpdateItem(ctx context.Context, userID uint, itemID uint, item eo.CartItem) error
	RemoveItem(ctx context.Context, userID uint, itemID uint) error
	Clear(ctx context.Context, userID uint) error
}

type ProductOptionRepository interface {
	List(ctx context.Context, productID *uint) ([]ec.ProductOption, error)
	Create(ctx context.Context, o *ec.ProductOption) error
	Update(ctx context.Context, id uint, o ec.ProductOption) error
	Delete(ctx context.Context, id uint) error
}

// Orders
type OrderRepository interface {
	CreateWithItems(ctx context.Context, o *eo.Order) error
	ListByUser(ctx context.Context, userID uint) ([]eo.Order, error)
	FindByID(ctx context.Context, id uint) (*eo.Order, error)
	FindWithItems(ctx context.Context, id uint) (*eo.Order, error)
	ListAll(ctx context.Context, status string) ([]eo.Order, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
	CountByStatus(ctx context.Context, status string) (int64, error)
	UpdatePaymentStatus(ctx context.Context, id uint, status string) error
}

// Repository is an aggregator passed into services
type Repository struct {
	Users          UserRepository
	Departments    DepartmentRepository
	Categories     CategoryRepository
	Products       ProductRepository
	ProductOptions ProductOptionRepository
	Orders         OrderRepository
	Carts          CartRepository
}
