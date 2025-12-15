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

	// ================= DB SETUP =================
	influxRepo, cleanupInflux := database.SetupInfluxDB(cfg)
	defer cleanupInflux()

	postgresRepo, cleanupPostgres := database.SetupPostgres(cfg)
	defer cleanupPostgres()

	cronSvc := service.NewCronService(influxRepo, postgresRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := utils.SetupSignalChannel()

	// ================= TIME ALIGNMENT =================
	now := time.Now().UTC()

	nextHour := now.Truncate(time.Hour).Add(time.Hour)

	nextDay := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0,
		time.UTC,
	).Add(24 * time.Hour)

	log.Printf("[INIT] Next hourly aggregation in %s\n", time.Until(nextHour).Round(time.Minute))
	log.Printf("[INIT] Next daily aggregation in %s\n", time.Until(nextDay).Round(time.Minute))

	// Timer awal untuk align jam
	hourlyTimer := time.NewTimer(time.Until(nextHour))

	// Reminder tiap 10 menit
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

	log.Println("[INIT] Cron service started")

	// ================= MAIN LOOP =================
	for {
		select {

		// ---------- REMINDER ----------
		case <-reminderTicker.C:
			now := time.Now().UTC()

			nextHourly := now.Truncate(time.Hour).Add(time.Hour)

			var nextDaily time.Time
			if now.Hour() < 24 {
				nextDaily = time.Date(
					now.Year(), now.Month(), now.Day(),
					0, 0, 0, 0,
					time.UTC,
				).Add(24 * time.Hour)
			}

			log.Printf(
				"[REMINDER] Hourly in %s | Daily in %s\n",
				time.Until(nextHourly).Round(time.Minute),
				time.Until(nextDaily).Round(time.Minute),
			)

		// ---------- HOURLY (+ DAILY) ----------
		case <-hourlyC:
			now := time.Now().UTC()

			// Jam yang BARU SELESAI
			targetHour := now.
				Add(-time.Hour).
				Truncate(time.Hour)

			log.Println("[RUN] HourlyAggregation for:", targetHour.Format(time.RFC3339))
			if _, err := cronSvc.HourlyAggregation(ctx, targetHour); err != nil {
				log.Println("[ERROR] HourlyAggregation:", err)
			}

			// ---------- DAILY ----------
			if targetHour.Hour() == 23 {
				targetDay := targetHour.Truncate(24 * time.Hour)

				log.Println("[RUN] DailyAggregation for:", targetDay.Format("2006-01-02"))
				if _, err := cronSvc.DailyAggregation(ctx, targetDay); err != nil {
					log.Println("[ERROR] DailyAggregation:", err)
				}
			}

			// Start ticker setelah alignment pertama
			if hourlyTicker == nil {
				hourlyTicker = time.NewTicker(time.Hour)
				hourlyC = hourlyTicker.C
			}

		// ---------- SHUTDOWN ----------
		case sig := <-sigChan:
			log.Println("[SHUTDOWN] Cron service stopped:", sig)
			return
		}
	}
}
