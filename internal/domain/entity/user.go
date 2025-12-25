package entity

import "metertronik/pkg/utils"

// type User struct {
//  ID        int64      `json:"id" gorm:"primaryKey"`
//  Email     string     `json:"email" gorm:"uniqueIndex;not null"`
//  Name      string     `json:"name" gorm:"not null"`
//  Password  string     `json:"-" gorm:"not null"`
//  Role      string     `json:"role" gorm:"default:user"`
//  IsActive  bool       `json:"is_active" gorm:"default:true"`
//  CreatedAt utils.TimeData  `json:"created_at" gorm:"autoCreateTime"`
//  UpdatedAt utils.TimeData  `json:"updated_at" gorm:"autoUpdateTime"`
// }

type User struct {
	ID        int64      `json:"id" gorm:"primaryKey"`
	Email     string     `json:"email" gorm:"uniqueIndex;not null"`
	Username  string     `json:"username" gorm:"not null"`
	Password  string     `json:"-" gorm:"not null"`
	Role      string     `json:"role" gorm:"default:user"`
	Status    string     `json:"status" gorm:"default:active"`
	Verified  bool       `json:"verified" gorm:"default:false"`
	CreatedAt utils.TimeData  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt utils.TimeData  `json:"updated_at" gorm:"autoUpdateTime"`
}

type Device struct {
	ID              int64          `json:"id" gorm:"primaryKey"`
	DeviceName      string         `json:"device_name" gorm:"not null"`
	DeviceType      string         `json:"device_type" gorm:"not null"`
	DeviceStatus    string         `json:"device_status" gorm:"not null"`
	DeviceLocation  string         `json:"device_location" gorm:"not null"`
	DeviceCreatedAt utils.TimeData `json:"device_created_at" gorm:"autoCreateTime"`
}
