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

	postgresRepo, _, cleanupPostgres := database.SetupPostgres(cfg)
	defer cleanupPostgres()

	cronSvc := service.NewCronService(influxRepo, postgresRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := utils.SetupSignalChannel()

	now := utils.TimeNow()

	nextHour := now.Truncate(time.Hour).Add(time.Hour)

	nextDay := now.StartOfDay().Add(24 * time.Hour)

	log.Printf("[INIT] Next hourly aggregation in %s\n", utils.TimeUntil(nextHour).Round(time.Minute))
	log.Printf("[INIT] Next daily aggregation in %s\n", utils.TimeUntil(nextDay).Round(time.Minute))

	hourlyTimer := time.NewTimer(utils.TimeUntil(nextHour))

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
			now := utils.TimeNow()

			devices, err := influxRepo.GetActiveDeviceIDs(ctx, 48)
			if err == nil && len(devices) > 0 {
				activeDevices = devices
				log.Printf("[REMINDER] Active devices updated: %v", activeDevices)
			}

			nextHourly := now.Truncate(time.Hour).Add(time.Hour)

			var nextDaily utils.TimeData
			if now.Time.Hour() < 24 {
				nextDaily = now.StartOfDay().Add(24 * time.Hour)
			}

			log.Printf(
				"[REMINDER] Hourly in %s | Daily in %s (UTC) | Active devices: %d\n",
				utils.TimeUntil(nextHourly).Round(time.Minute),
				utils.TimeUntil(nextDaily).Round(time.Minute),
				len(activeDevices),
			)

		case <-hourlyC:
			now := utils.TimeNow()

			devices, err := influxRepo.GetActiveDeviceIDs(ctx, 48)
			if err == nil && len(devices) > 0 {
				activeDevices = devices
			}

			targetHour := now.
				Add(-time.Hour).
				Truncate(time.Hour)

			log.Printf("[RUN] HourlyAggregation for (UTC): %s | Processing %d device(s)",
				targetHour.Format(), len(activeDevices))

			for _, deviceID := range activeDevices {
				log.Printf("[RUN] Processing device: %s", deviceID)
				if _, err := cronSvc.HourlyAggregation(ctx, targetHour.Time, deviceID); err != nil {
					log.Printf("[ERROR] HourlyAggregation for device %s: %v", deviceID, err)
				} else {
					log.Printf("[SUCCESS] HourlyAggregation completed for device: %s", deviceID)
				}
			}

			if targetHour.Time.Hour() == 23 {
				targetDay := targetHour.Truncate(24 * time.Hour)

				log.Printf("[RUN] DailyAggregation for: %s | Processing %d device(s)",
					targetDay.FormatLayout("2006-01-02"), len(activeDevices))

				for _, deviceID := range activeDevices {
					log.Printf("[RUN] Processing daily aggregation for device: %s", deviceID)
					if _, err := cronSvc.DailyAggregation(ctx, targetDay.Time, deviceID); err != nil {
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
