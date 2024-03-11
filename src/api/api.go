package api

import (
	"fmt"
	"time"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/helper"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/router"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	_ "github.com/Arshia-Izadyar/Fast-Gopher/src/docs"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/pkg/service_errors"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

func InitServer(cfg *config.Config) error {
	app := fiber.New(
		fiber.Config{
			JSONEncoder: sonic.Marshal,
			JSONDecoder: sonic.Unmarshal,
		},
	)

	swaggerInit(app)
	addMiddleware(app)
	registerRouters(app, cfg)

	err := app.Listen(fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		return err
	}
	return nil
}

func registerRouters(app *fiber.App, cfg *config.Config) {
	router.UserRouter(app, cfg)
	router.WhiteListAddRouter(app, cfg)
}

func addMiddleware(app *fiber.App) {
	app.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 60 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(helper.GenerateResponseWithError(&service_errors.ServiceError{EndUserMessage: "too many requests"}, false))
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		Storage:                nil,
		LimiterMiddleware:      limiter.SlidingWindow{},
		Duration:               0,
		Store:                  nil,
	}))

	app.Use(helmet.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}  \n",
	}))
}

func swaggerInit(app *fiber.App) {

	app.Get("/swagger/*", swagger.HandlerDefault) // default

}
