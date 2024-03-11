package router

import (
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/handler"
	authentication "github.com/Arshia-Izadyar/Fast-Gopher/src/api/middleware"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/gofiber/fiber/v2"
)

func UserRouter(app *fiber.App, cfg *config.Config) {
	h := handler.NewUserHandler(cfg)
	app.Post("/", h.TestHandler)
	app.Post("/login", h.LoginHandler)
	app.Get("/logout", authentication.New(cfg), h.Logout)
	app.Get("/google", h.GoogleLogin)
	app.Get("/auth/callback/google", h.LoginWithGoogleCode)
}
