package catalog

import (
	"context"

	ec "furniture-shop/internal/entities/catalog"
	"furniture-shop/internal/service"
	"furniture-shop/internal/storage"
)

type catalogService struct {
	departments storage.DepartmentRepository
	categories  storage.CategoryRepository
	products    storage.ProductRepository
}

func NewCatalogService(departments storage.DepartmentRepository, categories storage.CategoryRepository, products storage.ProductRepository) service.CatalogService {
	return &catalogService{departments: departments, categories: categories, products: products}
}

func (s *catalogService) ListDepartments(ctx context.Context) ([]ec.Department, error) {
	return s.departments.List(ctx)
}

func (s *catalogService) ListCategoriesByDepartment(ctx context.Context, departmentID uint) ([]ec.Category, error) {
	return s.categories.ListByDepartment(ctx, departmentID)
}

func (s *catalogService) ListProductsByCategory(ctx context.Context, categoryID uint) ([]ec.Product, error) {
	return s.products.ListByCategory(ctx, categoryID)
}

func (s *catalogService) GetProduct(ctx context.Context, id uint) (*ec.Product, error) {
	p, err := s.products.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = s.products.IncrementRecommendation(ctx, id)
	return p, nil
}

func (s *catalogService) SearchProducts(ctx context.Context, query string, limit int) ([]ec.Product, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.products.Search(ctx, query, limit)
}

func (s *catalogService) RecommendProducts(ctx context.Context, productID uint, limit int) ([]ec.Product, error) {
	p, err := s.products.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 4
	}
	return s.products.ListRecommendations(ctx, p, limit)
}
