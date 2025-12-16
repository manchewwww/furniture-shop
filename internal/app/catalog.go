package app

import (
    "context"

    "furniture-shop/internal/domain/repository"
    models "furniture-shop/internal/domain/entity"
)

type CatalogService interface {
    ListDepartments(ctx context.Context) ([]models.Department, error)
    ListCategoriesByDepartment(ctx context.Context, departmentID uint) ([]models.Category, error)
    ListProductsByCategory(ctx context.Context, categoryID uint) ([]models.Product, error)
    GetProduct(ctx context.Context, id uint) (*models.Product, error)
    SearchProducts(ctx context.Context, query string, limit int) ([]models.Product, error)
    RecommendProducts(ctx context.Context, productID uint, limit int) ([]models.Product, error)
}

type catalogService struct {
    departments repository.DepartmentRepository
    categories  repository.CategoryRepository
    products    repository.ProductRepository
}

func NewCatalogService(departments repository.DepartmentRepository, categories repository.CategoryRepository, products repository.ProductRepository) CatalogService {
    return &catalogService{departments: departments, categories: categories, products: products}
}

func (s *catalogService) ListDepartments(ctx context.Context) ([]models.Department, error) {
    return s.departments.List(ctx)
}

func (s *catalogService) ListCategoriesByDepartment(ctx context.Context, departmentID uint) ([]models.Category, error) {
    return s.categories.ListByDepartment(ctx, departmentID)
}

func (s *catalogService) ListProductsByCategory(ctx context.Context, categoryID uint) ([]models.Product, error) {
    return s.products.ListByCategory(ctx, categoryID)
}

func (s *catalogService) GetProduct(ctx context.Context, id uint) (*models.Product, error) {
    return s.products.FindByID(ctx, id)
}

func (s *catalogService) SearchProducts(ctx context.Context, query string, limit int) ([]models.Product, error) {
    if limit <= 0 { limit = 50 }
    return s.products.Search(ctx, query, limit)
}

func (s *catalogService) RecommendProducts(ctx context.Context, productID uint, limit int) ([]models.Product, error) {
    p, err := s.products.FindByID(ctx, productID)
    if err != nil { return nil, err }
    if limit <= 0 { limit = 4 }
    return s.products.ListRecommendations(ctx, p, limit)
}

