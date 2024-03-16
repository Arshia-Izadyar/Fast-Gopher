package router

import (
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/handler"
	authentication "github.com/Arshia-Izadyar/Fast-Gopher/src/api/middleware"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/gofiber/fiber/v2"
)

func WhiteListAddRouter(app fiber.Router, cfg *config.Config) {
	h := handler.NewWhiteListHandler(cfg)
	app.Get("/w", authentication.New(cfg), h.Add)
	app.Get("/rw", authentication.New(cfg), h.Remove)
}
