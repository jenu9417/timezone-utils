package main

import (
	"log"
	"timezone-utils/internal/api"
	"timezone-utils/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	router := api.SetupRouter(cfg)

	log.Printf("Server starting on %s...", cfg.Port)
	router.Run(":" + cfg.Port)
}
