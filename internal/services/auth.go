package services

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"

    "furniture-shop/internal/config"
    "furniture-shop/internal/database"
    "furniture-shop/internal/models"
)

func GenerateJWT(u *models.User, cfg *config.Config) (string, error) {
    claims := jwt.MapClaims{
        "sub": u.ID,
        "email": u.Email,
        "role": u.Role,
        "exp": time.Now().Add(24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(cfg.JWTSecret))
}

func Authenticate(email, password string) (*models.User, error) {
    var user models.User
    if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
        return nil, errors.New("невалиден имейл или парола")
    }
    if !user.CheckPassword(password) {
        return nil, errors.New("невалиден имейл или парола")
    }
    return &user, nil
}

