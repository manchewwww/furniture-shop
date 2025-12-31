package admin

import (
	"context"

	ec "furniture-shop/internal/entities/catalog"
	ei "furniture-shop/internal/entities/inventory"
	"furniture-shop/internal/service"
	"furniture-shop/internal/storage"
)

type adminService struct {
	depts   storage.DepartmentRepository
	cats    storage.CategoryRepository
	prods   storage.ProductRepository
	options storage.ProductOptionRepository
	stocks  storage.StockRepository
}

func NewAdminService(depts storage.DepartmentRepository, cats storage.CategoryRepository, prods storage.ProductRepository, options storage.ProductOptionRepository, stocks storage.StockRepository) service.AdminService {
	return &adminService{depts: depts, cats: cats, prods: prods, options: options, stocks: stocks}
}

func (s *adminService) ListDepartments(ctx context.Context) ([]ec.Department, error) {
	return s.depts.List(ctx)
}

func (s *adminService) CreateDepartment(ctx context.Context, d *ec.Department) error {
	return s.depts.Create(ctx, d)
}

func (s *adminService) UpdateDepartment(ctx context.Context, id uint, d ec.Department) error {
	return s.depts.Update(ctx, id, d)
}

func (s *adminService) DeleteDepartment(ctx context.Context, id uint) error {
	return s.depts.Delete(ctx, id)
}

func (s *adminService) ListCategories(ctx context.Context) ([]ec.Category, error) {
	return s.cats.ListAll(ctx)
}

func (s *adminService) CreateCategory(ctx context.Context, c *ec.Category) error {
	return s.cats.Create(ctx, c)
}

func (s *adminService) UpdateCategory(ctx context.Context, id uint, c ec.Category) error {
	return s.cats.Update(ctx, id, c)
}

func (s *adminService) DeleteCategory(ctx context.Context, id uint) error {
	return s.cats.Delete(ctx, id)
}

func (s *adminService) ListProducts(ctx context.Context) ([]ec.Product, error) {
	return s.prods.ListAll(ctx)
}

func (s *adminService) CreateProduct(ctx context.Context, p *ec.Product) error {
	return s.prods.Create(ctx, p)
}

func (s *adminService) UpdateProduct(ctx context.Context, id uint, p ec.Product) error {
	return s.prods.Update(ctx, id, p)
}

func (s *adminService) DeleteProduct(ctx context.Context, id uint) error {
	return s.prods.Delete(ctx, id)
}

func (s *adminService) ListProductOptions(ctx context.Context, productID *uint) ([]ec.ProductOption, error) {
	return s.options.List(ctx, productID)
}

func (s *adminService) CreateProductOption(ctx context.Context, o *ec.ProductOption) error {
	return s.options.Create(ctx, o)
}

func (s *adminService) UpdateProductOption(ctx context.Context, id uint, o ec.ProductOption) error {
	return s.options.Update(ctx, id, o)
}

func (s *adminService) DeleteProductOption(ctx context.Context, id uint) error {
	return s.options.Delete(ctx, id)
}

func (s *adminService) ListStock(ctx context.Context) ([]ei.Stock, error) {
	return s.stocks.List(ctx)
}

func (s *adminService) UpsertStock(ctx context.Context, material string, qty float64, unit string) error {
	return s.stocks.UpsertMaterial(ctx, material, qty, unit)
}
