package main

import (
	"log"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/postgres"
)

func main() {
	config.LoadConfig()
	config.LoadGoogleConfig()

	cfg := config.GetConfig()

	_, err := postgres.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	api.InitServer(cfg)

}
