package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"main/internal/model"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() *gorm.DB {
	godotenv.Load()

	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile("ca.pem")
	if err != nil {
		log.Fatal("Gagal membaca ca.pem:", err)
	}
	rootCertPool.AppendCertsFromPEM(pem)

	mysql.RegisterTLSConfig("custom-tls", &tls.Config{
		RootCAs: rootCertPool,
	})

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=custom-tls",
		os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	db, err := gorm.Open(gorm_mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}

	// Auto Migration
	err = db.AutoMigrate(
		&model.Consumer{},
		&model.ConsumerLimit{},
		&model.Transaction{},
	)
	if err != nil {
		log.Fatal("Gagal melakukan migration:", err)
	}

	log.Println("✓ Database connected and migrated successfully")
	return db
}
