package model

import "time"

type Consumer struct {
	ID        uint    `gorm:"primaryKey"`
	NIK       string  `gorm:"unique;not null;type:varchar(16)"`
	FullName  string  `gorm:"not null"`
	LegalName string  `gorm:"not null"`
	Salary    float64 `gorm:"type:decimal(15,2)"`
	CreatedAt time.Time
}

type ConsumerLimit struct {
	ID          uint    `gorm:"primaryKey"`
	ConsumerID  uint    `gorm:"index"`
	Tenor       int     `gorm:"not null"` // 1, 2, 3, 6 bulan
	LimitAmount float64 `gorm:"type:decimal(15,2)"`
}

type Transaction struct {
	ID                uint    `gorm:"primaryKey"`
	ConsumerID        uint    `gorm:"index"`
	ContractNumber    string  `gorm:"unique;not null"`
	OTR               float64 `gorm:"type:decimal(15,2)"`
	AdminFee          float64 `gorm:"type:decimal(15,2)"`
	InstallmentAmount float64 `gorm:"type:decimal(15,2)"`
	InterestAmount    float64 `gorm:"type:decimal(15,2)"`
	AssetName         string  `gorm:"type:varchar(255)"`
}
