package main

import (
	"log"
	"metertronik/pkg/config"
	"metertronik/pkg/database"
	"metertronik/internal/service"
	"context"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	influxRepo, cleanupInflux := database.SetupInfluxDB(cfg)
	defer cleanupInflux()

	scv := service.NewCronService(influxRepo)

	scv.HourlyAggregation(context.Background())
}