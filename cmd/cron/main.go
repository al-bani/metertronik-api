package main

import (
	"context"
	"log"

	"time"

	"metertronik/internal/service"
	"metertronik/pkg/config"
	"metertronik/pkg/database"
	"metertronik/pkg/utils"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	influxRepo, cleanupInflux := database.SetupInfluxDB(cfg)
	defer cleanupInflux()

	postgresRepo, cleanupPostgres := database.SetupPostgres(cfg)
	defer cleanupPostgres()

	scv := service.NewCronService(influxRepo, postgresRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := utils.SetupSignalChannel()

	hourlyTicker := time.NewTicker(1 * time.Hour)
	defer hourlyTicker.Stop()

	dailyTicker := time.NewTicker(24 * time.Hour)
	defer dailyTicker.Stop()

	log.Println("Cron service started")

	for {
		select {
		case <-hourlyTicker.C:
			log.Println("Starting HourlyAggregation")
			scv.HourlyAggregation(ctx)

		case <-dailyTicker.C:
			log.Println("Starting DailyAggregation")
			scv.DailyAggregation(ctx)

		case sig := <-sigChan:
			log.Println("Closed...", sig)
			return
		}
	}
}
