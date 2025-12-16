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

	cronSvc := service.NewCronService(influxRepo, postgresRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := utils.SetupSignalChannel()

	now := time.Now().In(utils.WIBLocation())

	nextHour := now.Truncate(time.Hour).Add(time.Hour)

	nextDay := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0,
		utils.WIBLocation(),
	).Add(24 * time.Hour)

	log.Printf("[INIT] Next hourly aggregation in %s\n", time.Until(nextHour).Round(time.Minute))
	log.Printf("[INIT] Next daily aggregation in %s\n", time.Until(nextDay).Round(time.Minute))

	hourlyTimer := time.NewTimer(time.Until(nextHour))

	reminderTicker := time.NewTicker(10 * time.Minute)

	var hourlyTicker *time.Ticker
	hourlyC := hourlyTimer.C

	defer func() {
		hourlyTimer.Stop()
		reminderTicker.Stop()
		if hourlyTicker != nil {
			hourlyTicker.Stop()
		}
	}()

	activeDevices, err := influxRepo.GetActiveDeviceIDs(ctx, 48)
	if err != nil {
		log.Printf("[WARNING] Failed to get active device IDs: %v", err)
		activeDevices = []string{"device-001"}
	}
	log.Printf("[INIT] Active devices found: %v", activeDevices)
	log.Println("[INIT] Cron service started")

	for {
		select {

		case <-reminderTicker.C:
			now := time.Now().In(utils.WIBLocation())

			devices, err := influxRepo.GetActiveDeviceIDs(ctx, 48)
			if err == nil && len(devices) > 0 {
				activeDevices = devices
				log.Printf("[REMINDER] Active devices updated: %v", activeDevices)
			}

			nextHourly := now.Truncate(time.Hour).Add(time.Hour)

			var nextDaily time.Time
			if now.Hour() < 24 {
				nextDaily = time.Date(
					now.Year(), now.Month(), now.Day(),
					0, 0, 0, 0,
					utils.WIBLocation(),
				).Add(24 * time.Hour)
			}

			log.Printf(
				"[REMINDER] Hourly in %s | Daily in %s (WIB) | Active devices: %d\n",
				time.Until(nextHourly).Round(time.Minute),
				time.Until(nextDaily).Round(time.Minute),
				len(activeDevices),
			)

		case <-hourlyC:
			now := time.Now().In(utils.WIBLocation())

			devices, err := influxRepo.GetActiveDeviceIDs(ctx, 48)
			if err == nil && len(devices) > 0 {
				activeDevices = devices
			}

			targetHour := now.
				Add(-time.Hour).
				Truncate(time.Hour)

			log.Printf("[RUN] HourlyAggregation for (WIB): %s | Processing %d device(s)",
				targetHour.Format(time.RFC3339), len(activeDevices))

			for _, deviceID := range activeDevices {
				log.Printf("[RUN] Processing device: %s", deviceID)
				if _, err := cronSvc.HourlyAggregation(ctx, targetHour, deviceID); err != nil {
					log.Printf("[ERROR] HourlyAggregation for device %s: %v", deviceID, err)
				} else {
					log.Printf("[SUCCESS] HourlyAggregation completed for device: %s", deviceID)
				}
			}

			if targetHour.Hour() == 23 {
				targetDay := targetHour.Truncate(24 * time.Hour)

				log.Printf("[RUN] DailyAggregation for: %s | Processing %d device(s)",
					targetDay.Format("2006-01-02"), len(activeDevices))

				for _, deviceID := range activeDevices {
					log.Printf("[RUN] Processing daily aggregation for device: %s", deviceID)
					if _, err := cronSvc.DailyAggregation(ctx, targetDay, deviceID); err != nil {
						log.Printf("[ERROR] DailyAggregation for device %s: %v", deviceID, err)
					} else {
						log.Printf("[SUCCESS] DailyAggregation completed for device: %s", deviceID)
					}
				}
			}

			if hourlyTicker == nil {
				hourlyTicker = time.NewTicker(time.Hour)
				hourlyC = hourlyTicker.C
			}

		case sig := <-sigChan:
			log.Println("[SHUTDOWN] Cron service stopped:", sig)
			return
		}
	}
}
