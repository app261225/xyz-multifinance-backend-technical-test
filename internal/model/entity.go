package model

import "time"

// Consumer represents a customer of PT XYZ Multifinance
type Consumer struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	NIK            string          `gorm:"unique;not null;type:varchar(16)" json:"nik"`
	FullName       string          `gorm:"not null;type:varchar(255)" json:"full_name"`
	LegalName      string          `gorm:"not null;type:varchar(255)" json:"legal_name"`
	PlaceOfBirth   string          `gorm:"type:varchar(255)" json:"place_of_birth"`
	DateOfBirth    time.Time       `json:"date_of_birth"`
	Salary         float64         `gorm:"type:decimal(15,2)" json:"salary"`
	KTPPhoto       string          `gorm:"type:text" json:"ktp_photo"`
	SelfiePhoto    string          `gorm:"type:text" json:"selfie_photo"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      *time.Time      `gorm:"index" json:"deleted_at,omitempty"`
	ConsumerLimits []ConsumerLimit `gorm:"foreignKey:ConsumerID" json:"limits,omitempty"`
	Transactions   []Transaction   `gorm:"foreignKey:ConsumerID" json:"transactions,omitempty"`
}

// ConsumerLimit represents the credit limit for a consumer
type ConsumerLimit struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ConsumerID  uint      `gorm:"index;not null" json:"consumer_id"`
	Consumer    Consumer  `json:"consumer,omitempty"`
	Tenor       int       `gorm:"not null" json:"tenor"` // 1, 2, 3, 6 bulan
	LimitAmount float64   `gorm:"type:decimal(15,2);not null" json:"limit_amount"`
	UsedAmount  float64   `gorm:"type:decimal(15,2);default:0" json:"used_amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Transaction represents a financial transaction
type Transaction struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	ConsumerID        uint      `gorm:"index;not null" json:"consumer_id"`
	Consumer          Consumer  `json:"consumer,omitempty"`
	ContractNumber    string    `gorm:"unique;not null;type:varchar(255)" json:"contract_number"`
	Tenor             int       `gorm:"not null" json:"tenor"` // 1, 2, 3, 6 bulan
	OTR               float64   `gorm:"type:decimal(15,2);not null" json:"otr"`
	AdminFee          float64   `gorm:"type:decimal(15,2)" json:"admin_fee"`
	InstallmentAmount float64   `gorm:"type:decimal(15,2);not null" json:"installment_amount"`
	InterestAmount    float64   `gorm:"type:decimal(15,2)" json:"interest_amount"`
	AssetName         string    `gorm:"type:varchar(255)" json:"asset_name"`
	Status            string    `gorm:"type:varchar(50);default:'ACTIVE'" json:"status"` // ACTIVE, COMPLETED, DEFAULTED
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
