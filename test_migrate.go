package main

import (
	"bizkit-backend/config"
	"bizkit-backend/internal/model"
	"fmt"
)

func main() {
	config.ConnectDB()
	err := config.DB.AutoMigrate(&model.ReceivablePayment{}, &model.ReceivablePaymentItem{})
	if err != nil {
		fmt.Println("MIGRATE ERROR:", err)
	} else {
		fmt.Println("MIGRATE SUCCESS")
	}
}
