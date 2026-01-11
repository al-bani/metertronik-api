package service

import (
	"context"
	"log"
	"math"
	"metertronik/internal/domain/entity"
	"metertronik/internal/domain/repository"
	"metertronik/pkg/utils"
	"strings"
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
	log.Printf("Checking Data...")
	previousData, err := s.RedisRealtimeRepo.GetLatestElectricity(ctx, data.DeviceID)

	if err != nil {
		data.PowerSurge = 0
		data.PSPercent = 0
	} else if previousData == nil {
		data.PowerSurge = 0
		data.PSPercent = 0
	} else {
		log.Printf("Using previousData")
		data.PowerSurge = math.Abs(data.Power - previousData.Power)
		// PSPercent harus dihitung terhadap nilai sebelumnya (selama bukan 0),
		// jangan diblok dengan minBaseline besar (ini bikin trigger threshold sering miss).
		if previousData.Power != 0 {
			data.PSPercent = math.Abs((data.PowerSurge / previousData.Power) * 100)
		} else {
			data.PSPercent = 0
		}
	}

	if data.CreatedAt.Time.IsZero() {
		data.CreatedAt = utils.TimeNow()
	}

	errInflux := s.influxRepo.SaveRealTimeElectricity(ctx, data)
	_ = errInflux

	if previousData == nil {
		log.Printf("Save Latest Data to Redis")
		if err := s.RedisRealtimeRepo.SetLatestElectricity(ctx, data.DeviceID, data); err != nil {
			_ = err
		}

		if err := s.RedisRealtimeRepo.SaveElectricityHistory(ctx, data.DeviceID, data, 5*time.Minute); err != nil {
			_ = err
		}
		return nil
	}

	proximityValue := ProximityValue(previousData, data)
	if !proximityValue {
		log.Printf("no data exceeds the threshold")
		return nil
	}

	threshold := 5.0
	var reasons []string
	if data.PowerSurge > 100.0 {
		reasons = append(reasons, "PowerSurge>100")
	}
	if data.PSPercent > 5.0 {
		reasons = append(reasons, "PSPercent>5")
	}
	if diff := percentageDiff(data.Power, previousData.Power); diff >= threshold {
		reasons = append(reasons, "Power>=5%")
	}
	if diff := percentageDiff(data.Voltage, previousData.Voltage); diff >= threshold {
		reasons = append(reasons, "Voltage>=5%")
	}
	if diff := percentageDiff(data.Current, previousData.Current); diff >= threshold {
		reasons = append(reasons, "Current>=5%")
	}
	if diff := percentageDiff(data.Energy, previousData.Energy); diff >= threshold {
		reasons = append(reasons, "Energy>=5%")
	}
	if diff := percentageDiff(data.PowerFactor, previousData.PowerFactor); diff >= threshold {
		reasons = append(reasons, "PowerFactor>=5%")
	}
	if diff := percentageDiff(data.Frequency, previousData.Frequency); diff >= threshold {
		reasons = append(reasons, "Frequency>=5%")
	}
	if len(reasons) > 0 {
		log.Printf("data exceeds the threshold: %s", strings.Join(reasons, ", "))
	} else {
		log.Printf("data exceeds the threshold")
	}
	log.Printf("Save Latest Data to Redis")
	if err := s.RedisRealtimeRepo.SetLatestElectricity(ctx, data.DeviceID, data); err != nil {
		_ = err
	}

	if err := s.RedisRealtimeRepo.SaveElectricityHistory(ctx, data.DeviceID, data, 5*time.Minute); err != nil {
		_ = err
	}

	return nil
}

func percentageDiff(current, previous float64) float64 {
	if previous == 0 {
		// kalau previous = 0 dan current != 0, ini perubahan signifikan (hindari miss update)
		if current == 0 {
			return 0
		}
		return math.Inf(1)
	}
	return math.Abs(((current - previous) / previous) * 100)
}

func ProximityValue(previousData *entity.RealTimeElectricity, data *entity.RealTimeElectricity) bool {
	if previousData == nil || data == nil {
		return false
	}

	threshold := 5.0

	diffPower := percentageDiff(data.Power, previousData.Power)
	diffVoltage := percentageDiff(data.Voltage, previousData.Voltage)
	diffCurrent := percentageDiff(data.Current, previousData.Current)
	diffEnergy := percentageDiff(data.Energy, previousData.Energy)
	diffPF := percentageDiff(data.PowerFactor, previousData.PowerFactor)
	diffFreq := percentageDiff(data.Frequency, previousData.Frequency)

	if data.PowerSurge > 100.0 || data.PSPercent > 5.0 {
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
