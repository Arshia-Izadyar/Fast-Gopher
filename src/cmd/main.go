package main

import (
	"log"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/router"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/postgres"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	cfg := config.GetConfig()
	router.UserRouter(app, cfg)

	_, err := postgres.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app.Listen(":8000")
}
