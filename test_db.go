package main

import (
	"bizkit-backend/config"
	"bizkit-backend/internal/repository"
	"fmt"
)

func main() {
	config.ConnectDB()

	// Cek Sale 28
	sale, err := repository.GetSaleByID(28)
	if err != nil {
		fmt.Println("Error Sale 28:", err)
	} else {
		fmt.Printf("Sale 28 OK: %#v\n", sale.InvoiceNumber)
	}

	// Cek Endpoint Payments
	items, err := repository.GetPaymentsBySaleID(28)
	if err != nil {
		fmt.Println("Error Payments 28:", err)
	} else {
		fmt.Printf("Payments 28 OK: %d items\n", len(items))
	}
}
