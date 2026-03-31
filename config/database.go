package config

import (
	"fmt"
	"log"
	"os"

	"bizkit-backend/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	db.AutoMigrate(
		&model.Role{},
		&model.User{},
		&model.Category{},
		&model.Brand{},
		&model.Unit{},
		&model.VariantCategory{},
		&model.VariantOption{},
		&model.Product{},
		&model.PaymentMethod{},
		&model.Promo{},
		&model.Sale{},
		&model.SaleItem{},
		&model.SaleItemVariant{},
		&model.Attendance{},
		&model.Shift{},
		&model.Setting{},
		&model.Outlet{},
		&model.PriceCategory{},
		&model.ProductPrice{},
		&model.PromoItem{},
		&model.PromoSpecialPrice{},
		&model.PromoVoucher{},
		&model.Customer{},
		&model.ReceivablePayment{},
		&model.ReceivablePaymentItem{},
	)

	DB = db
	log.Println("Database berhasil terhubung!")
}
