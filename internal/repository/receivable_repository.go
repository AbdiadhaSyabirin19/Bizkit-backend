package repository

import (
	"bizkit-backend/config"
	"bizkit-backend/internal/model"
)

func CreateReceivablePayment(payment *model.ReceivablePayment) error {
	return config.DB.Create(payment).Error
}

func UpdateReceivablePayment(payment *model.ReceivablePayment) error {
	return config.DB.Save(payment).Error
}

func DeleteReceivablePayment(id uint) error {
	return config.DB.Delete(&model.ReceivablePayment{}, id).Error
}

func GetReceivablePaymentByID(id uint) (*model.ReceivablePayment, error) {
	var payment model.ReceivablePayment
	err := config.DB.Preload("Items").Preload("Items.Sale").Preload("PaymentMethod").First(&payment, id).Error

	// Break json loop
	for i := range payment.Items {
		payment.Items[i].ReceivablePayment = nil
	}

	return &payment, err
}

func GetPaymentsBySaleID(saleID uint) ([]model.ReceivablePaymentItem, error) {
	var items []model.ReceivablePaymentItem
	err := config.DB.Preload("ReceivablePayment").
		Preload("ReceivablePayment.PaymentMethod").
		Preload("ReceivablePayment.User").
		Where("sale_id = ?", saleID).
		Find(&items).Error

	// Break json loop
	for i := range items {
		if items[i].ReceivablePayment != nil {
			items[i].ReceivablePayment.Items = nil
		}
	}

	return items, err
}

func FixMigrate() error {
	return config.DB.AutoMigrate(&model.ReceivablePayment{}, &model.ReceivablePaymentItem{})
}

func SyncSales() error {
	var sales []model.Sale
	config.DB.Where("payment_status = ?", "lunas").Find(&sales)

	for _, sale := range sales {
		if sale.AmountPaid == 0 && sale.GrandTotal > 0 {
			sale.AmountPaid = sale.GrandTotal
			config.DB.Save(&sale)
		}

		var count int64
		config.DB.Model(&model.ReceivablePaymentItem{}).Where("sale_id = ?", sale.ID).Count(&count)
		if count == 0 {
			pmID := sale.PaymentMethodID
			if pmID == 0 {
				pmID = 1
			}

			payment := model.ReceivablePayment{
				PaymentDate:     sale.CreatedAt,
				CustomerName:    sale.CustomerName,
				PaymentMethodID: pmID,
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
			config.DB.Create(&payment)
		}
	}
	return nil
}

func GetUnpaidSales(customerName string) ([]model.Sale, error) {
	var sales []model.Sale
	query := config.DB.Where("payment_status != ?", "lunas")
	if customerName != "" {
		query = query.Where("customer_name = ?", customerName)
	}
	err := query.Find(&sales).Error
	return sales, err
}

func RecalculateSalePayment(saleID uint) error {
	var sale model.Sale
	if err := config.DB.First(&sale, saleID).Error; err != nil {
		return err
	}

	var totalPaid float64
	config.DB.Model(&model.ReceivablePaymentItem{}).
		Where("sale_id = ?", saleID).
		Select("COALESCE(SUM(amount_paid), 0)").
		Scan(&totalPaid)

	status := "piutang"
	if totalPaid >= sale.GrandTotal {
		status = "lunas"
	}

	return config.DB.Model(&sale).Updates(map[string]interface{}{
		"amount_paid":    totalPaid,
		"payment_status": status,
	}).Error
}
