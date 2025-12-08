package entity

import (
	"metertronik/pkg/utils"
)

type RealTimeElectricity struct {
	ID          int64          `json:"id" gorm:"primaryKey"`
	DeviceID    string         `json:"device_id" gorm:"index;not null"`
	Voltage     float64        `json:"voltage" gorm:"type:decimal(10,2);not null"`
	Current     float64        `json:"current" gorm:"type:decimal(10,3);not null"`
	Power       float64        `json:"power" gorm:"type:decimal(10,2);not null"`
	TotalEnergy float64        `json:"total_energy" gorm:"type:decimal(15,3);not null"`
	PowerFactor float64        `json:"power_factor" gorm:"type:decimal(4,2);not null"`
	Frequency   float64        `json:"frequency" gorm:"type:decimal(5,2);not null"`
	CreatedAt   utils.TimeData `json:"created_at" gorm:"autoCreateTime"`
}

type HourlyElectricity struct {
	ID           int64          `json:"id" gorm:"primaryKey"`
	DeviceID     string         `json:"device_id" gorm:"index;not null"`
	UsageKWh     float64        `json:"usage_kwh" gorm:"type:decimal(10,3);not null"`
	TotalCost    float64        `json:"total_cost" gorm:"type:decimal(15,2);not null"`
	AvgVoltage   float64        `json:"avg_voltage" gorm:"type:decimal(10,2)"`
	AvgCurrent   float64        `json:"avg_current" gorm:"type:decimal(10,3)"`
	AvgPower     float64        `json:"avg_power" gorm:"type:decimal(10,2)"`
	AvgFrequency float64        `json:"avg_frequency" gorm:"type:decimal(5,2)"`
	MinPower     float64        `json:"min_power" gorm:"type:decimal(10,2)"`
	MaxPower     float64        `json:"max_power" gorm:"type:decimal(10,2)"`
	CreatedAt    utils.TimeData `json:"created_at" gorm:"autoCreateTime"`
}

type DailyElectricity struct {
	ID         int64          `json:"id" gorm:"primaryKey"`
	DeviceID   string         `json:"device_id" gorm:"index;not null"`
	Date       utils.TimeData `json:"date" gorm:"type:date;not null"`
	UsageKWh   float64        `json:"usage_kwh" gorm:"type:decimal(10,3);not null"`
	TotalCost  float64        `json:"total_cost" gorm:"type:decimal(15,2);not null"`
	AvgVoltage float64        `json:"avg_voltage"`
	MinPower   float64        `json:"min_power"`
	MaxPower   float64        `json:"max_power"`
	CreatedAt  utils.TimeData `json:"created_at" gorm:"autoCreateTime"`
}

type MonthlyElectricity struct {
	ID        int64          `json:"id" gorm:"primaryKey"`
	DeviceID  string         `json:"device_id" gorm:"index;not null"`
	MonthDate utils.TimeData `json:"month_date" gorm:"type:date;not null"`
	UsageKWh  float64        `json:"usage_kwh"`
	TotalCost float64        `json:"total_cost"`
	CreatedAt utils.TimeData `json:"created_at" gorm:"autoCreateTime"`
}
