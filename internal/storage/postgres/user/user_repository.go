package user

import (
    "context"
    "gorm.io/gorm"

    eu "furniture-shop/internal/entities/user"
    "furniture-shop/internal/storage"
)

type UserRepository struct { db *gorm.DB }

func NewUserRepository(db *gorm.DB) storage.UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) Create(ctx context.Context, u *eu.User) error {
    return r.db.WithContext(ctx).Create(u).Error
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*eu.User, error) {
    var u eu.User
    if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil { return nil, err }
    return &u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*eu.User, error) {
    var u eu.User
    if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil { return nil, err }
    return &u, nil
}

