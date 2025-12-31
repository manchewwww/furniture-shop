package postgres

import (
	"gorm.io/gorm"

	"furniture-shop/internal/storage"
	pgadmin "furniture-shop/internal/storage/postgres/catalog"
	pginv "furniture-shop/internal/storage/postgres/inventory"
	pgorders "furniture-shop/internal/storage/postgres/orders"
	pguser "furniture-shop/internal/storage/postgres/user"
)

// NewRepository wires Postgres-backed repositories
func NewRepository(db *gorm.DB) *storage.Repository {
	return &storage.Repository{
		Users:          pguser.NewUserRepository(db),
		Departments:    pgadmin.NewDepartmentRepository(db),
		Categories:     pgadmin.NewCategoryRepository(db),
		Products:       pgadmin.NewProductRepository(db),
		ProductOptions: pgadmin.NewProductOptionRepository(db),
		Orders:         pgorders.NewOrderRepository(db),
		Carts:          pgorders.NewCartRepository(db),
		Stock:          pginv.NewStockRepository(db),
	}
}
