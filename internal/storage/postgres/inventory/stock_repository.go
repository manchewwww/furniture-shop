package inventory

import (
	"context"

	"gorm.io/gorm"

	ei "furniture-shop/internal/entities/inventory"
	"furniture-shop/internal/storage"
)

type StockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) storage.StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) FindByMaterial(ctx context.Context, material string) (float64, error) {
	var s ei.Stock
	if err := r.db.WithContext(ctx).Where("material_name = ?", material).First(&s).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	return s.QuantityAvailable, nil
}

func (r *StockRepository) UpsertMaterial(ctx context.Context, material string, qty float64, unit string) error {
	var s ei.Stock
	if err := r.db.WithContext(ctx).Where("material_name = ?", material).First(&s).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			s = ei.Stock{MaterialName: material, QuantityAvailable: qty, Unit: unit}
			return r.db.WithContext(ctx).Create(&s).Error
		}
		return err
	}
	s.QuantityAvailable = qty
	s.Unit = unit
	return r.db.WithContext(ctx).Save(&s).Error
}

func (r *StockRepository) AdjustQuantity(ctx context.Context, material string, delta float64) error {
	var s ei.Stock
	if err := r.db.WithContext(ctx).Where("material_name = ?", material).First(&s).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			s = ei.Stock{MaterialName: material, QuantityAvailable: delta, Unit: "pcs"}
			return r.db.WithContext(ctx).Create(&s).Error
		}
		return err
	}
	s.QuantityAvailable += delta
	return r.db.WithContext(ctx).Save(&s).Error
}

func (r *StockRepository) List(ctx context.Context) ([]ei.Stock, error) {
	var items []ei.Stock
	if err := r.db.WithContext(ctx).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
