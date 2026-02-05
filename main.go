package main

import (
	"fmt"
	"log"
	"time"

	"main/config"
	"main/internal/model"
	"main/internal/repository"
	"main/internal/usecase" // Tambah import ini
)

func main() {
	// 1. Koneksi Database
	db := config.ConnectDB()

	// 2. Setup Layer (Merakit Aplikasi)
	// Layer bawah: Repository (butuh DB)
	consumerRepo := repository.NewConsumerRepository(db)

	// Layer tengah: Usecase (butuh Repository)
	consumerService := usecase.NewConsumerUsecase(consumerRepo)

	// 3. Siapkan Data Dummy
	// Kita coba data yang valid
	newConsumer := model.Consumer{
		NIK:       fmt.Sprintf("NIK-%d", time.Now().Unix()), // Unik
		FullName:  "Citra Usecase Tester",
		LegalName: "Citra Resmi",
		Salary:    25000000,
		CreatedAt: time.Now(),
	}

	fmt.Println("Mencoba mendaftarkan konsumen via Usecase...")

	// 4. Panggil Usecase (Bukan Repository langsung)
	err := consumerService.RegisterConsumer(&newConsumer)
	if err != nil {
		log.Fatal("Gagal register:", err)
	}

	fmt.Println("Berhasil: Data masuk melalui Validasi Usecase!")
}
