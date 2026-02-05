package usecase

import (
	"errors"
	"log"

	// Sesuaikan nama module dengan go.mod kamu
	"main/internal/model"
	"main/internal/repository"
)

// 1. Interface: Kontrak layanan apa yang tersedia untuk konsumen
type ConsumerUsecase interface {
	RegisterConsumer(consumer *model.Consumer) error
}

// 2. Struct: Usecase butuh Repository untuk akses DB
type consumerUsecase struct {
	repo repository.ConsumerRepository
}

// 3. Constructor: Merakit Usecase dengan Repository yang dibutuhkan
func NewConsumerUsecase(repo repository.ConsumerRepository) ConsumerUsecase {
	return &consumerUsecase{repo}
}

// 4. Method: Logika Bisnis (Validasi & Flow)
func (u *consumerUsecase) RegisterConsumer(consumer *model.Consumer) error {
	// --- LOGIKA BISNIS DI SINI ---

	// Validasi 1: Cek apakah data kosong
	if consumer.FullName == "" || consumer.NIK == "" {
		return errors.New("nama dan NIK tidak boleh kosong")
	}

	// Validasi 2: Cek Gaji (Syarat Multifinance)
	if consumer.Salary < 0 {
		return errors.New("gaji tidak valid (negatif)")
	}

	// Jika semua validasi lolos, baru panggil Repository untuk simpan
	log.Println("Logika Bisnis OK. Menyimpan ke Repository...")
	return u.repo.Create(consumer)
}
