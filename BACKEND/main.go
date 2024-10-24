package main

import (
	"exchangeapp/config"
	"exchangeapp/router"
)

func main() {
	config.InitConfig()
	port := config.AppConfig.App.Port
	if port == "" {
		port = ":8080"
	}
	r := router.SetupRouter()
	r.Run(port)
}