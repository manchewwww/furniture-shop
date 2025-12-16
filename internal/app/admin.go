package app

import (
    "context"
    "furniture-shop/internal/domain/repository"
    models "furniture-shop/internal/domain/entity"
)

type AdminService interface {
    // Departments
    ListDepartments(ctx context.Context) ([]models.Department, error)
    CreateDepartment(ctx context.Context, d *models.Department) error
    UpdateDepartment(ctx context.Context, id uint, d models.Department) error
    DeleteDepartment(ctx context.Context, id uint) error
    // Categories
    ListCategories(ctx context.Context) ([]models.Category, error)
    CreateCategory(ctx context.Context, c *models.Category) error
    UpdateCategory(ctx context.Context, id uint, c models.Category) error
    DeleteCategory(ctx context.Context, id uint) error
    // Products
    ListProducts(ctx context.Context) ([]models.Product, error)
    CreateProduct(ctx context.Context, p *models.Product) error
    UpdateProduct(ctx context.Context, id uint, p models.Product) error
    DeleteProduct(ctx context.Context, id uint) error
    // Product Options
    ListProductOptions(ctx context.Context, productID *uint) ([]models.ProductOption, error)
    CreateProductOption(ctx context.Context, o *models.ProductOption) error
    UpdateProductOption(ctx context.Context, id uint, o models.ProductOption) error
    DeleteProductOption(ctx context.Context, id uint) error
}

type adminService struct {
    depts   repository.DepartmentRepository
    cats    repository.CategoryRepository
    prods   repository.ProductRepository
    options repository.ProductOptionRepository
}

func NewAdminService(depts repository.DepartmentRepository, cats repository.CategoryRepository, prods repository.ProductRepository, options repository.ProductOptionRepository) AdminService {
    return &adminService{depts: depts, cats: cats, prods: prods, options: options}
}

func (s *adminService) ListDepartments(ctx context.Context) ([]models.Department, error) { return s.depts.List(ctx) }
func (s *adminService) CreateDepartment(ctx context.Context, d *models.Department) error { return s.depts.Create(ctx, d) }
func (s *adminService) UpdateDepartment(ctx context.Context, id uint, d models.Department) error { return s.depts.Update(ctx, id, d) }
func (s *adminService) DeleteDepartment(ctx context.Context, id uint) error { return s.depts.Delete(ctx, id) }

func (s *adminService) ListCategories(ctx context.Context) ([]models.Category, error) { return s.cats.ListAll(ctx) }
func (s *adminService) CreateCategory(ctx context.Context, c *models.Category) error { return s.cats.Create(ctx, c) }
func (s *adminService) UpdateCategory(ctx context.Context, id uint, c models.Category) error { return s.cats.Update(ctx, id, c) }
func (s *adminService) DeleteCategory(ctx context.Context, id uint) error { return s.cats.Delete(ctx, id) }

func (s *adminService) ListProducts(ctx context.Context) ([]models.Product, error) { return s.prods.ListAll(ctx) }
func (s *adminService) CreateProduct(ctx context.Context, p *models.Product) error { return s.prods.Create(ctx, p) }
func (s *adminService) UpdateProduct(ctx context.Context, id uint, p models.Product) error { return s.prods.Update(ctx, id, p) }
func (s *adminService) DeleteProduct(ctx context.Context, id uint) error { return s.prods.Delete(ctx, id) }

func (s *adminService) ListProductOptions(ctx context.Context, productID *uint) ([]models.ProductOption, error) { return s.options.List(ctx, productID) }
func (s *adminService) CreateProductOption(ctx context.Context, o *models.ProductOption) error { return s.options.Create(ctx, o) }
func (s *adminService) UpdateProductOption(ctx context.Context, id uint, o models.ProductOption) error { return s.options.Update(ctx, id, o) }
func (s *adminService) DeleteProductOption(ctx context.Context, id uint) error { return s.options.Delete(ctx, id) }

