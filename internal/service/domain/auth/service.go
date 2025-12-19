package auth

import (
    "context"
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v4"

    eu "furniture-shop/internal/entities/user"
    "furniture-shop/internal/service"
    "furniture-shop/internal/storage"
)

type authService struct {
    users     storage.UserRepository
    jwtSecret string
}

func NewAuthService(users storage.UserRepository, jwtSecret string) service.AuthService {
    return &authService{users: users, jwtSecret: jwtSecret}
}

func (s *authService) GenerateJWT(u *eu.User) (string, error) {
    claims := jwt.MapClaims{
        "sub":   u.ID,
        "email": u.Email,
        "role":  u.Role,
        "exp":   time.Now().Add(24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.jwtSecret))
}

func (s *authService) Authenticate(ctx context.Context, email, password string) (*eu.User, error) {
    user, err := s.users.FindByEmail(ctx, email)
    if err != nil { return nil, errors.New("invalid email or password") }
    if !user.CheckPassword(password) { return nil, errors.New("invalid email or password") }
    return user, nil
}

func (s *authService) CreateUser(ctx context.Context, u *eu.User) error { return s.users.Create(ctx, u) }

