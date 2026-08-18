package main

import (
	"log"

	"github.com/keep/sunny/config"
	"github.com/keep/sunny/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
