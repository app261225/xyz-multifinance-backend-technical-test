package main

import (
	"crypto/tls" // Tambahkan ini
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	godotenv.Load()

	// 1. Baca file CA Certificate
	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile("ca.pem")
	if err != nil {
		log.Fatal("Gagal membaca file ca.pem:", err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		log.Fatal("Gagal melampirkan sertifikat CA")
	}

	// 2. Registrasi TLS Config yang BENAR
	// Gunakan &tls.Config, bukan &mysql.Config
	err = mysql.RegisterTLSConfig("custom-tls", &tls.Config{
		RootCAs: rootCertPool,
	})
	if err != nil {
		log.Fatal("Gagal registrasi TLS:", err)
	}

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=custom-tls",
		user, pass, host, port, name)

	db, err := gorm.Open(gorm_mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal koneksi GORM: %v", err)
	}

	return db
}

func main() {
	db := ConnectDB()
	fmt.Println("✅ Koneksi Berhasil dengan SSL/TLS Aman!")

	// Cek koneksi fisik
	sqlDB, _ := db.DB()
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Ping gagal:", err)
	}
}
