package router

import (
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/handler"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/middleware"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/middleware/authentication"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/gofiber/fiber/v2"
)

func WhiteListAddRouter(app fiber.Router, cfg *config.Config) {
	h := handler.NewWhiteListHandler(cfg)
	app.Get("/w", h.Add)
	app.Get("/w/premium", authentication.New(cfg), middleware.Premium(), h.PremiumAdd)
	app.Get("/rw", authentication.New(cfg), h.Remove)
}
