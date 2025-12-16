package gormrepo

import (
    "context"

    "gorm.io/gorm"

    "furniture-shop/internal/domain/repository"
    models "furniture-shop/internal/domain/entity"
)

type ProductRepository struct { db *gorm.DB }

func NewProductRepository(db *gorm.DB) repository.ProductRepository { return &ProductRepository{db: db} }

func (r *ProductRepository) ListByCategory(ctx context.Context, categoryID uint) ([]models.Product, error) {
    var out []models.Product
    if err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Find(&out).Error; err != nil {
        return nil, err
    }
    return out, nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id uint) (*models.Product, error) {
    var p models.Product
    if err := r.db.WithContext(ctx).Preload("Options").First(&p, id).Error; err != nil {
        return nil, err
    }
    return &p, nil
}

func (r *ProductRepository) Search(ctx context.Context, query string, limit int) ([]models.Product, error) {
    var items []models.Product
    like := "%" + query + "%"
    if err := r.db.WithContext(ctx).
        Where("name ILIKE ? OR short_description ILIKE ?", like, like).
        Limit(limit).
        Find(&items).Error; err != nil {
        return nil, err
    }
    return items, nil
}

func (r *ProductRepository) ListRecommendations(ctx context.Context, p *models.Product, limit int) ([]models.Product, error) {
    var rec []models.Product
    if err := r.db.WithContext(ctx).
        Where("category_id = ? AND id <> ?", p.CategoryID, p.ID).
        Limit(limit).
        Find(&rec).Error; err != nil {
        return nil, err
    }
    return rec, nil
}

func (r *ProductRepository) ListAll(ctx context.Context) ([]models.Product, error) {
    var items []models.Product
    if err := r.db.WithContext(ctx).Find(&items).Error; err != nil { return nil, err }
    return items, nil
}

func (r *ProductRepository) Create(ctx context.Context, p *models.Product) error {
    return r.db.WithContext(ctx).Create(p).Error
}

func (r *ProductRepository) Update(ctx context.Context, id uint, p models.Product) error {
    p.ID = 0
    return r.db.WithContext(ctx).Model(&models.Product{}).Where("id = ?", id).Updates(p).Error
}

func (r *ProductRepository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}
