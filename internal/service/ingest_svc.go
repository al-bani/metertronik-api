package service

import (
	"context"
	"log"
	"math"
	"metertronik/internal/domain/entity"
	"metertronik/internal/domain/repository"
	"metertronik/pkg/utils"
	"time"
)

type IngestService struct {
	influxRepo        repository.InfluxRepo
	RedisRealtimeRepo repository.RedisRealtimeRepo
}

func NewIngestService(influxRepo repository.InfluxRepo, RedisRealtimeRepo repository.RedisRealtimeRepo) *IngestService {
	return &IngestService{
		influxRepo:        influxRepo,
		RedisRealtimeRepo: RedisRealtimeRepo,
	}
}

func (s *IngestService) ProcessRealTimeElectricity(ctx context.Context, data *entity.RealTimeElectricity) error {
	log.Printf("\n\nProcessing electricity data for device: %s", data.DeviceID)

	previousData, err := s.RedisRealtimeRepo.GetLatestElectricity(ctx, data.DeviceID)

	if err != nil {
		log.Printf("Error getting previous electricity data: %v", err)
		data.PowerSurge = 0
		data.PSPercent = 0
	} else if previousData == nil {
		data.PowerSurge = 0
		data.PSPercent = 0
	} else {
		data.PowerSurge = math.Abs(data.Power - previousData.Power)
		minBaseline := 50.0

		if previousData.Power >= minBaseline {
			data.PSPercent = math.Abs((data.PowerSurge / previousData.Power) * 100)
		} else {
			log.Printf("Previous power is 0, setting PSPercent to 0")
			data.PSPercent = 0
		}
	}

	if data.CreatedAt.Time.IsZero() {
		data.CreatedAt = utils.TimeNow()
	}

	errInflux := s.influxRepo.SaveRealTimeElectricity(ctx, data)

	if errInflux != nil {
		log.Printf("Error saving real time electricity to influx: %v", errInflux)
	} else {
		log.Println("Saving data to influxDB : ", data)
	}

	// Jika previousData == nil, ini adalah data pertama, selalu cache
	if previousData == nil {
		log.Printf("First data for device %s, caching immediately", data.DeviceID)
		if err := s.RedisRealtimeRepo.SetLatestElectricity(ctx, data.DeviceID, data); err != nil {
			log.Printf("Failed saving latest cache: %v", err)
		} else {
			log.Println("Updated latest cache data")
		}

		if err := s.RedisRealtimeRepo.SaveElectricityHistory(ctx, data.DeviceID, data, 5*time.Minute); err != nil {
			log.Printf("Failed saving history cache: %v", err)
		}
		return nil
	}

	changed, _, err := s.RedisRealtimeRepo.HasChanged(ctx, data.DeviceID, data)
	if err != nil {
		log.Printf("Error comparing cache: %v", err)
	}

	if !changed {
		log.Printf("No change for device %s (skip caching)", data.DeviceID)
		return nil
	}

	proximityValue := ProximityValue(previousData, data)
	if !proximityValue {
		log.Printf("No significant change for device %s, skipping caching", data.DeviceID)
		return nil
	}

	if err := s.RedisRealtimeRepo.SetLatestElectricity(ctx, data.DeviceID, data); err != nil {
		log.Printf("Failed saving latest cache: %v", err)
	} else {
		log.Println("Updated latest cache data")
	}

	if err := s.RedisRealtimeRepo.SaveElectricityHistory(ctx, data.DeviceID, data, 5*time.Minute); err != nil {
		log.Printf("Failed saving history cache: %v", err)
	}

	return nil
}

func percentageDiff(current, previous float64) float64 {
	if previous == 0 {
		return 0
	}
	return math.Abs(((current - previous) / previous) * 100)
}

func ProximityValue(previousData *entity.RealTimeElectricity, data *entity.RealTimeElectricity) bool {
	if previousData == nil || data == nil {
		return false
	}

	threshold := 10.0

	diffPower := percentageDiff(data.Power, previousData.Power)
	diffVoltage := percentageDiff(data.Voltage, previousData.Voltage)
	diffCurrent := percentageDiff(data.Current, previousData.Current)
	diffEnergy := percentageDiff(data.Energy, previousData.Energy)
	diffPF := percentageDiff(data.PowerFactor, previousData.PowerFactor)
	diffFreq := percentageDiff(data.Frequency, previousData.Frequency)

	if data.PowerSurge > 500.0 || data.PSPercent > 15.0 {
		return true
	}

	if diffPower >= threshold ||
		diffVoltage >= threshold ||
		diffCurrent >= threshold ||
		diffEnergy >= threshold ||
		diffPF >= threshold ||
		diffFreq >= threshold {
		return true
	}

	return false
}
