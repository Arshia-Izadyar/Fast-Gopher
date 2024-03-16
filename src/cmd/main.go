package main

import (
	"log"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/cmd/cmd"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/cache"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/postgres"
)

// @title Internal auth
// @version 0.1
// @description internal service for Auth
// @termsOfService Kir
// @contact.name API Support
// @contact.email arshiaa104@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:4000
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name authorization
var W *cmd.WorkerPool

func main() {
	config.LoadConfig()
	config.LoadGoogleConfig()
	cfg := config.GetConfig()
	err := cache.InitRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}

	_, err = postgres.ConnectDB(cfg)
	defer postgres.CloseDB()
	defer cache.Close()

	if err != nil {
		log.Fatal(err)
	}
	W = cmd.New(300)
	api.InitServer(cfg)

}
