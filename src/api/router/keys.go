package router

import (
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/handler"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/gofiber/fiber/v2"
)

func UserRouter(app fiber.Router, cfg *config.Config) {
	h := handler.NewKeyHandler(cfg)
	// app.Get("/logout", authentication.New(cfg), h.Logout)
	app.Post("/refresh", h.Refresh)
	app.Post("/key", h.GenerateKey)
	app.Post("/key/tk", h.GenerateTokenFromKey)
}
