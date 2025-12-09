package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"metertronik/internal/service"
	"metertronik/pkg/config"
	"metertronik/pkg/database"
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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	log.Println("🚀 Cron scheduler started, running HourlyAggregation for the first time...")
	runHourlyAggregation(ctx, scv)

	log.Println("⏰ Scheduler active, will run HourlyAggregation every hour")

	for {
		select {
		case <-ticker.C:
			log.Println("⏰ Time to run HourlyAggregation")
			runHourlyAggregation(ctx, scv)
		case sig := <-sigChan:
			log.Printf("📴 Received signal %v, stopping scheduler...", sig)
			return
		}
	}
}

func runHourlyAggregation(ctx context.Context, scv *service.CronService) {
	_, err := scv.HourlyAggregation(ctx)
	if err != nil {
		log.Printf("❌ Error running HourlyAggregation: %v", err)
	} else {
		log.Println("✅ HourlyAggregation executed successfully")
	}
}
