package catalog

import (
	"context"

	"gorm.io/gorm"

	ec "furniture-shop/internal/entities/catalog"
	"furniture-shop/internal/storage"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) storage.ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) ListByCategory(ctx context.Context, categoryID uint) ([]ec.Product, error) {
	var out []ec.Product
	if err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id uint) (*ec.Product, error) {
	var p ec.Product
	if err := r.db.WithContext(ctx).Preload("Options").First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Search(ctx context.Context, query string, limit int) ([]ec.Product, error) {
	var items []ec.Product
	like := "%" + query + "%"
	if err := r.db.WithContext(ctx).
		Where("name ILIKE ? OR short_description ILIKE ?", like, like).
		Limit(limit).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ProductRepository) ListRecommendations(ctx context.Context, p *ec.Product, limit int) ([]ec.Product, error) {
	var rec []ec.Product
	q := r.db.WithContext(ctx).Model(&ec.Product{}).
		Joins("LEFT JOIN recommendation_counters rc ON rc.product_id = products.id").
		Where("products.category_id = ? AND products.id <> ?", p.CategoryID, p.ID).
		Order("COALESCE(rc.count,0) DESC, products.id ASC").
		Limit(limit)
	if err := q.Find(&rec).Error; err != nil {
		return nil, err
	}
	return rec, nil
}

func (r *ProductRepository) ListAll(ctx context.Context) ([]ec.Product, error) {
	var items []ec.Product
	if err := r.db.WithContext(ctx).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ProductRepository) Create(ctx context.Context, p *ec.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *ProductRepository) Update(ctx context.Context, id uint, p ec.Product) error {
	return r.db.WithContext(ctx).Model(&ec.Product{}).Where("id = ?", id).
		Select("name", "short_description", "long_description", "base_price", "base_production_time_days", "category_id", "image_url",
			"default_width", "default_height", "default_depth", "base_material").
		Updates(p).Error
}

func (r *ProductRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&ec.Product{}, id).Error
}

func (r *ProductRepository) IncrementRecommendation(ctx context.Context, productID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var rc ec.RecommendationCounter
		if err := tx.Where("product_id = ?", productID).First(&rc).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				rc = ec.RecommendationCounter{ProductID: productID, Count: 1}
				return tx.Create(&rc).Error
			}
			return err
		}
		return tx.Model(&rc).UpdateColumn("count", gorm.Expr("count + 1")).Error
	})
}
