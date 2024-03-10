package api

import (
	"fmt"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/router"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func InitServer(cfg *config.Config) error {
	app := fiber.New(
		fiber.Config{
			JSONEncoder: sonic.Marshal,
			JSONDecoder: sonic.Unmarshal,
		},
	)

	addMiddleware(app, cfg)
	registerRouters(app, cfg)

	err := app.Listen(fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		return err
	}
	return nil
}

func registerRouters(app *fiber.App, cfg *config.Config) {
	router.UserRouter(app, cfg)

}

func addMiddleware(app *fiber.App, cfg *config.Config) {
	// limiterConfig := limiter.Config{

	// 	Max:        2,                // max count of connections
	// 	Expiration: 30 * time.Second, // expiration time of the limit

	// 	LimitReached: func(c *fiber.Ctx) error {
	// 		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"change": "this"}) // called when a request hits the limit
	// 	},
	// }
	fmt.Println("here")
	app.Use(helmet.New())
	// app.Use(limiter.New(limiterConfig))
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}  \n",
	}))
}
