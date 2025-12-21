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
	Energy      float64        `json:"energy" gorm:"type:decimal(15,3);not null"`
	PowerFactor float64        `json:"power_factor" gorm:"type:decimal(4,2);not null"`
	Frequency   float64        `json:"frequency" gorm:"type:decimal(5,2);not null"`
	PowerSurge  float64        `json:"power_surge" gorm:"type:decimal(10,2);not null"`
	PSPercent   float64        `json:"power_surge_percentage" gorm:"type:decimal(4,2);not null"`
	CreatedAt   utils.TimeData `json:"created_at" gorm:"autoCreateTime"`
}

type HourlyElectricity struct {
	DeviceID   string  `json:"device_id" gorm:"column:device_id;type:varchar(50);not null"`
	Energy     float64 `json:"energy" gorm:"column:energy;type:decimal(10,3);not null"`
	TotalCost  float64 `json:"total_cost" gorm:"column:total_cost;type:decimal(15,2);not null"`
	AvgVoltage float64 `json:"avg_voltage" gorm:"column:avg_voltage;type:decimal(10,2)"`
	AvgCurrent float64 `json:"avg_current" gorm:"column:avg_current;type:decimal(10,3)"`
	AvgPower   float64 `json:"avg_power" gorm:"column:avg_power;type:decimal(10,2)"`
	MinPower   float64 `json:"min_power" gorm:"column:min_power;type:decimal(10,2)"`
	MaxPower   float64 `json:"max_power" gorm:"column:max_power;type:decimal(10,2)"`

	TS        utils.TimeData `json:"ts" gorm:"column:ts;type:timestamptz;not null"`
	CreatedAt utils.TimeData `json:"created_at" gorm:"autoCreateTime"`
}

type DailyElectricity struct {
	DeviceID string `json:"device_id" gorm:"column:device_id;type:varchar(50);not null"`

	Energy     float64 `json:"energy" gorm:"column:energy;type:decimal(10,3);not null"`
	TotalCost  float64 `json:"total_cost" gorm:"column:total_cost;type:decimal(15,2);not null"`
	AvgVoltage float64 `json:"avg_voltage" gorm:"column:avg_voltage;type:decimal(10,2)"`
	AvgCurrent float64 `json:"avg_current" gorm:"column:avg_current;type:decimal(10,3)"`
	AvgPower   float64 `json:"avg_power" gorm:"column:avg_power;type:decimal(10,2)"`
	MinPower   float64 `json:"min_power" gorm:"column:min_power;type:decimal(10,2)"`
	MaxPower   float64 `json:"max_power" gorm:"column:max_power;type:decimal(10,2)"`

	Day       utils.TimeData `json:"day" gorm:"column:day;type:date;not null"`
	CreatedAt utils.TimeData `json:"created_at" gorm:"autoCreateTime"`
}

type MonthlyElectricity struct {
	DeviceID string `json:"device_id" gorm:"column:device_id;type:varchar(50);not null"`
	Month    utils.TimeData `json:"month" gorm:"column:month;type:date;not null"`

	Energy    float64        `json:"energy" gorm:"column:energy;type:decimal(10,3);not null"`
	TotalCost float64        `json:"total_cost" gorm:"column:total_cost;type:decimal(15,2);not null"`
	CreatedAt utils.TimeData `json:"created_at" gorm:"autoCreateTime"`
}
