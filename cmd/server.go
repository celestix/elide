package cmd

import (
	"elide/internal/api"
	"elide/internal/api/config"
	"log"
)

func runServer() {
	defer log.Println("[Elide][MAIN]: Started")
	config.Debug = true
	engine := api.CreateEngine(protocol, port, appId, apiHash, botToken)
	engine.Run()
}
