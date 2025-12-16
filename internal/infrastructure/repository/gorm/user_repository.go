package gormrepo

import (
    "context"

    "gorm.io/gorm"

    "furniture-shop/internal/domain/repository"
    "furniture-shop/internal/models"
)

type UserRepository struct { db *gorm.DB }

func NewUserRepository(db *gorm.DB) repository.UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
    return r.db.WithContext(ctx).Create(u).Error
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
    var u models.User
    if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
    var u models.User
    if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
        return nil, err
    }
    return &u, nil
}

