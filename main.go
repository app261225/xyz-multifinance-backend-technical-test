package main

import (
	"fmt"
	"log"

	"main/config"
	"main/internal/model"
)

func main() {
	// 1. Panggil koneksi dari folder config
	db := config.ConnectDB()
	fmt.Println("Terhubung ke Database")

	// 2. Jalankan Migrasi (Membuat tabel di Aiven Cloud)
	fmt.Println("Melakukan migrasi tabel...")
	err := db.AutoMigrate(&model.Consumer{}, &model.ConsumerLimit{}, &model.Transaction{})
	if err != nil {
		log.Fatal("❌ Gagal Migrasi:", err)
	}

	fmt.Println("Berhasil: Tabel otomatis dibuat di Cloud!")
}
