package http

import (
	"log"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"furniture-shop/internal/config"
	"furniture-shop/internal/service"
)

type Server struct {
	app *fiber.App
	svc *service.Service
}

func NewServer(svc *service.Service) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if verrs, ok := err.(validator.ValidationErrors); ok {
				out := make([]fiber.Map, 0, len(verrs))
				for _, fe := range verrs {
					out = append(out, fiber.Map{
						"field": fe.Field(),
						"tag":   fe.Tag(),
						"param": fe.Param(),
					})
				}
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "invalid input",
					"errors":  out,
				})
			}
			if fe, ok := err.(*fiber.Error); ok {
				return c.Status(fe.Code).JSON(fiber.Map{"message": fe.Message})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "internal server error"})
		},
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(config.Configurations.CORSOrigins, ","),
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PATCH,DELETE,PUT",
		AllowHeaders:     "Authorization,Content-Type",
	}))
	app.Static("/uploads", "./uploads")
	s := &Server{app: app, svc: svc}
	buildRoutes(s)
	return s
}

func (s *Server) Run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server listening on :%s", port)
	return s.app.Listen(":" + port)
}
