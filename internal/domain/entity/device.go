package entity

import "metertronik/pkg/utils"

type Device struct {
	ID int64 `json:"id" gorm:"primaryKey"`
	DeviceID string `json:"device_id" gorm:"not null"`
	DeviceName string `json:"device_name" gorm:"not null"`
	DeviceSecret string `json:"device_secret" gorm:"not null"`
	Paired bool `json:"paired" gorm:"default:false"`
	PairedAt utils.TimeData `json:"paired_at" gorm:"autoCreateTime"`
	CreatedAt utils.TimeData `json:"created_at" gorm:"autoCreateTime"`
	LastSeen utils.TimeData `json:"last_seen" gorm:"autoUpdateTime"`
}

type DeviceUser struct {
	ID int64 `json:"id" gorm:"primaryKey"`
	DeviceID int64 `json:"device_id" gorm:"not null"`
	UserID int64 `json:"user_id" gorm:"not null"`
	Role string `json:"role" gorm:"default:owner"`
	CreatedAt utils.TimeData `json:"created_at" gorm:"autoCreateTime"`
}