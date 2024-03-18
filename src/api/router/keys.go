package router

import (
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/handler"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/middleware/authentication"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/gofiber/fiber/v2"
)

func UserRouter(app fiber.Router, cfg *config.Config) {
	h := handler.NewKeyHandler(cfg)
	// app.Get("/logout", authentication.New(cfg), h.Logout)
	app.Post("/auth/refresh", h.Refresh)
	app.Post("/key", h.GenerateKey)                                    // 1. first the install guaranteed call
	app.Post("/key/tk", h.GenerateTokenFromKey)                        // 2. if they enter a new key for their device
	app.Get("/show", authentication.New(cfg), h.ShowActiveSessions)    // show all of active and white listed devices
	app.Delete("/rm", authentication.New(cfg), h.RemoveDevice)         // remove a device
	app.Delete("/rm/all", authentication.New(cfg), h.RemoveAllDevices) // remove all devices except this one
}
