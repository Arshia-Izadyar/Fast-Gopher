package api

import (
	"fmt"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/router"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	_ "github.com/Arshia-Izadyar/Fast-Gopher/src/docs"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// "github.com/gofiber/fiber/v2/middleware/limiter"
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
	api := app.Group("/api")
	v1 := api.Group("/v1")
	router.UserRouter(v1, cfg)
	router.WhiteListAddRouter(v1, cfg)
}

func addMiddleware(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// app.Use(limiter.New(limiter.Config{
	// 	Max:        20,
	// 	Expiration: 60 * time.Second,
	// 	LimitReached: func(c *fiber.Ctx) error {
	// 		return c.Status(fiber.StatusTooManyRequests).JSON(helper.GenerateResponseWithError(&service_errors.ServiceError{EndUserMessage: "too many requests"}, false))
	// 	},
	// 	SkipFailedRequests:     false,
	// 	SkipSuccessfulRequests: false,
	// 	Storage:                nil,
	// 	LimiterMiddleware:      limiter.SlidingWindow{},
	// 	Duration:               0,
	// 	Store:                  nil,
	// }))

	// app.Use(helmet.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${latency} ${method} ${path} \n",
	}))
}

func swaggerInit(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault) // default
}
