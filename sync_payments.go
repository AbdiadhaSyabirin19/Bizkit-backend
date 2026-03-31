package main

import (
	"bizkit-backend/config"
	"bizkit-backend/internal/model"
	"fmt"
)

func main() {
	config.ConnectDB()

	var sales []model.Sale
	// Find all LUNAS sales
	config.DB.Where("payment_status = ?", "lunas").Find(&sales)

	for _, sale := range sales {
		if sale.AmountPaid == 0 && sale.GrandTotal > 0 {
			sale.AmountPaid = sale.GrandTotal
			config.DB.Save(&sale)
		}

		// Check if payment history exists
		var count int64
		config.DB.Model(&model.ReceivablePaymentItem{}).Where("sale_id = ?", sale.ID).Count(&count)
		if count == 0 {
			// Backfill payment
			payment := model.ReceivablePayment{
				PaymentDate:     sale.CreatedAt,
				CustomerName:    sale.CustomerName,
				PaymentMethodID: sale.PaymentMethodID,
				Notes:           "Pembayaran otomatis (Sinkronisasi sistem lama)",
				Amount:          sale.GrandTotal,
				UserID:          sale.UserID,
				Items: []model.ReceivablePaymentItem{
					{
						SaleID:     sale.ID,
						AmountPaid: sale.GrandTotal,
					},
				},
			}
			if sale.PaymentMethodID == 0 {
				payment.PaymentMethodID = 1 // default cache
			}
			err := config.DB.Create(&payment).Error
			if err != nil {
				fmt.Println("Gagal sinkron Sale", sale.ID, ":", err)
			} else {
				fmt.Println("Sukses sinkron Sale", sale.ID)
			}
		}
	}
	fmt.Println("Sinkronisasi Selesai!")
}
