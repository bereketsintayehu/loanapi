package main

import (
	"loan/config"
	"loan/config/db"
	"loan/delivery/routers"
)

func main() {
	config.InitiEnvConfigs()
	db.ConnectDB(config.EnvConfigs.MongoURI)

	router := routers.SetupRouter()

	router.Run(config.EnvConfigs.ServerPort)
}
