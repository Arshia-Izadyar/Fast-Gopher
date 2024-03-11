package main

import (
	"log"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/postgres"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/redis"
)

func main() {
	config.LoadConfig()
	config.LoadGoogleConfig()
	cfg := config.GetConfig()

	err := redis.InitRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}

	_, err = postgres.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	api.InitServer(cfg)

}
