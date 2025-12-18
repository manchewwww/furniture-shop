package middleware

import (
    "os"

    "github.com/gofiber/fiber/v2"
    jwtware "github.com/gofiber/jwt/v3"
    "github.com/golang-jwt/jwt/v5"
)

func JWTAuth() fiber.Handler {
    jwtSecret := os.Getenv("JWT_SECRET")
    return jwtware.New(jwtware.Config{
        SigningKey: []byte(jwtSecret),
        ContextKey: "jwt",
        SuccessHandler: func(c *fiber.Ctx) error {
            user := c.Locals("jwt").(*jwt.Token)
            claims := user.Claims.(jwt.MapClaims)
            if v, ok := claims["sub"].(float64); ok { c.Locals("user_id", uint(v)) }
            if v, ok := claims["email"].(string); ok { c.Locals("user_email", v) }
            if v, ok := claims["role"].(string); ok { c.Locals("user_role", v) }
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

