package entity

import (
	"metertronik/pkg/utils"
)

type Tarrifs struct {
    ID				uint64     `gorm:"primaryKey;column:id"`
    TypeTarrif		string     `gorm:"column:type_tarrif;type:varchar(20);not null"`
    PowerVA			int        `gorm:"column:power_va;not null"`
    PricePerKwh		float64    `gorm:"column:price_per_kwh"`
    EffectiveFrom	utils.TimeData `gorm:"column:effective_from;not null"`
    EffectiveTo		utils.TimeData `gorm:"column:effective_to"`  
    CreatedAt		utils.TimeData  `gorm:"column:created_at;autoCreateTime"`
}