package repository

import (
	"main/internal/model"

	"gorm.io/gorm"
)

// 1. Interface: Kontrak yang bisa dilakukan repository
type ConsumerRepository interface {
	Create(consumer *model.Consumer) error
}

// 2. Struct: Implementasi nyata yang memegang koneksi DB
type consumerRepository struct {
	db *gorm.DB
}

// 3. Constructor: Fungsi untuk membuat instance repository baru
func NewConsumerRepository(db *gorm.DB) ConsumerRepository {
	return &consumerRepository{db}
}

// 4. Method: Fungsi Create/Simpan data
func (r *consumerRepository) Create(consumer *model.Consumer) error {
	// GORM akan otomatis generate SQL INSERT INTO consumers ...
	return r.db.Create(consumer).Error
}
