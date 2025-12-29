package middleware

import (
	"furniture-shop/internal/config"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

func JWTAuth() fiber.Handler {
	jwtSecret := config.Env.JWTSecret
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(jwtSecret),
		ContextKey: "jwt",
		SuccessHandler: func(c *fiber.Ctx) error {
			token, ok := c.Locals("jwt").(*jwt.Token)
			if !ok || token == nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid token"})
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid claims"})
			}
			if v, ok := claims["sub"].(float64); ok {
				c.Locals("user_id", uint(v))
			}
			if v, ok := claims["email"].(string); ok {
				c.Locals("user_email", v)
			}
			if v, ok := claims["role"].(string); ok {
				c.Locals("user_role", v)
			}
			return c.Next()
		},
	})
}

func RequireAdmin(c *fiber.Ctx) error {
	if c.Locals("user_role") != "admin" {
		return c.Status(403).JSON(fiber.Map{"message": "forbidden"})
	}
	return c.Next()
}
