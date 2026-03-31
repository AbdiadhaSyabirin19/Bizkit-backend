package main

import (
	"bizkit-backend/config"
	"bizkit-backend/internal/model"
	"bizkit-backend/internal/service"
	"fmt"
)

func main() {
	config.ConnectDB()

	// Get first user
	var user model.User
	if err := config.DB.First(&user).Error; err != nil {
		fmt.Println("No user found")
		return
	}

	// Get first payment method
	var pm model.PaymentMethod
	if err := config.DB.First(&pm).Error; err != nil {
		fmt.Println("No payment method found")
		return
	}

	// Get first product
	var product model.Product
	if err := config.DB.First(&product).Error; err != nil {
		fmt.Println("No product found")
		return
	}

	req := service.SaleRequest{
		PaymentMethodID: pm.ID,
		CustomerName:    "Test Customer",
		Source:          "dashboard",
		Items: []service.SaleItemRequest{
			{
				ProductID: product.ID,
				Quantity:  1,
			},
		},
	}

	sale, err := service.CreateSale(req, user.ID)
	if err != nil {
		fmt.Println("CreateSale failed:", err)
		return
	}

	fmt.Println("CreateSale SUCCESS! Invoice:", sale.InvoiceNumber)
}
