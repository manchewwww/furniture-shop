package services

import (
    "context"
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"

    "furniture-shop/internal/domain/repository"
    "furniture-shop/internal/models"
)

type AuthService interface {
    GenerateJWT(u *models.User) (string, error)
    Authenticate(ctx context.Context, email, password string) (*models.User, error)
    CreateUser(ctx context.Context, u *models.User) error
}

type authService struct {
    users     repository.UserRepository
    jwtSecret string
}

func NewAuthService(users repository.UserRepository, jwtSecret string) AuthService {
    return &authService{users: users, jwtSecret: jwtSecret}
}

func (s *authService) GenerateJWT(u *models.User) (string, error) {
    claims := jwt.MapClaims{
        "sub": u.ID,
        "email": u.Email,
        "role": u.Role,
        "exp": time.Now().Add(24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.jwtSecret))
}

func (s *authService) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
    user, err := s.users.FindByEmail(ctx, email)
    if err != nil {
        return nil, errors.New("invalid email or password")
    }
    if !user.CheckPassword(password) {
        return nil, errors.New("invalid email or password")
    }
    return user, nil
}

func (s *authService) CreateUser(ctx context.Context, u *models.User) error {
    return s.users.Create(ctx, u)
}
