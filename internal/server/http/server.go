package http

import (
	"log"
	"os"
	"strings"

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
	app := fiber.New()
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
