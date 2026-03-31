package service

import (
	"errors"
	"time"

	"bizkit-backend/config"
	"bizkit-backend/internal/model"
	"bizkit-backend/internal/repository"
)

type ReceivablePaymentItemRequest struct {
	SaleID     uint    `json:"sale_id" binding:"required"`
	AmountPaid float64 `json:"amount_paid" binding:"required"`
}

type ReceivablePaymentRequest struct {
	PaymentDate     string                         `json:"payment_date" binding:"required"`
	CustomerName    string                         `json:"customer_name"`
	PaymentMethodID uint                           `json:"payment_method_id" binding:"required"`
	Notes           string                         `json:"notes"`
	Amount          float64                        `json:"amount" binding:"required"`
	Items           []ReceivablePaymentItemRequest `json:"items" binding:"required"`
}

func CreateReceivablePayment(req ReceivablePaymentRequest, userID uint) (*model.ReceivablePayment, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("item pembayaran piutang tidak boleh kosong")
	}

	date, err := time.Parse(time.RFC3339, req.PaymentDate)
	if err != nil {
		date, err = time.Parse("2006-01-02T15:04:05.000Z", req.PaymentDate) // legacy format
		if err != nil {
			date, err = time.Parse("2006-01-02", req.PaymentDate)
			if err != nil {
				return nil, errors.New("format tanggal tidak valid")
			}
		}
	}

	tx := config.DB.Begin()

	var paymentItems []model.ReceivablePaymentItem
	for _, itemReq := range req.Items {
		sale, err := repository.GetSaleByID(itemReq.SaleID)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("transaksi penjualan tidak ditemukan")
		}

		sisa := sale.GrandTotal - sale.AmountPaid
		if itemReq.AmountPaid > sisa {
			tx.Rollback()
			return nil, errors.New("jumlah bayar melebihi sisa tunggakan")
		}

		paymentItems = append(paymentItems, model.ReceivablePaymentItem{
			SaleID:     itemReq.SaleID,
			AmountPaid: itemReq.AmountPaid,
		})
	}

	payment := model.ReceivablePayment{
		PaymentDate:     date,
		CustomerName:    req.CustomerName,
		PaymentMethodID: req.PaymentMethodID,
		Notes:           req.Notes,
		Amount:          req.Amount,
		UserID:          userID,
		Items:           paymentItems,
	}

	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal menyimpan pembayaran")
	}

	for _, item := range paymentItems {
		var totalPaid float64
		tx.Model(&model.ReceivablePaymentItem{}).
			Where("sale_id = ?", item.SaleID).
			Select("COALESCE(SUM(amount_paid), 0)").
			Scan(&totalPaid)

		var sale model.Sale
		tx.First(&sale, item.SaleID)

		status := "piutang"
		if totalPaid >= sale.GrandTotal {
			status = "lunas"
		}

		tx.Model(&sale).Updates(map[string]interface{}{
			"amount_paid":    totalPaid,
			"payment_status": status,
		})
	}

	tx.Commit()

	return repository.GetReceivablePaymentByID(payment.ID)
}

func UpdateReceivablePayment(id uint, req ReceivablePaymentRequest) (*model.ReceivablePayment, error) {
	existing, err := repository.GetReceivablePaymentByID(id)
	if err != nil {
		return nil, errors.New("pembayaran tidak ditemukan")
	}

	date, err := time.Parse(time.RFC3339, req.PaymentDate)
	if err != nil {
		date, err = time.Parse("2006-01-02T15:04:05.000Z", req.PaymentDate)
		if err != nil {
			date, err = time.Parse("2006-01-02", req.PaymentDate)
			if err != nil {
				return nil, errors.New("format tanggal tidak valid")
			}
		}
	}

	tx := config.DB.Begin()

	var oldSaleIDs []uint
	for _, it := range existing.Items {
		oldSaleIDs = append(oldSaleIDs, it.SaleID)
	}

	tx.Where("receivable_payment_id = ?", id).Delete(&model.ReceivablePaymentItem{})

	var paymentItems []model.ReceivablePaymentItem
	for _, itemReq := range req.Items {
		sale, err := repository.GetSaleByID(itemReq.SaleID)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("transaksi penjualan tidak ditemukan")
		}

		var paidByOthers float64
		tx.Model(&model.ReceivablePaymentItem{}).
			Where("sale_id = ? AND receivable_payment_id != ?", itemReq.SaleID, id).
			Select("COALESCE(SUM(amount_paid), 0)").
			Scan(&paidByOthers)

		sisa := sale.GrandTotal - paidByOthers
		if itemReq.AmountPaid > sisa {
			tx.Rollback()
			return nil, errors.New("jumlah bayar melebihi sisa tunggakan")
		}

		paymentItems = append(paymentItems, model.ReceivablePaymentItem{
			ReceivablePaymentID: id,
			SaleID:              itemReq.SaleID,
			AmountPaid:          itemReq.AmountPaid,
		})
	}

	for _, it := range paymentItems {
		if err := tx.Create(&it).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("gagal menyimpan item pembayaran")
		}
	}

	existing.PaymentDate = date
	existing.CustomerName = req.CustomerName
	existing.PaymentMethodID = req.PaymentMethodID
	existing.Notes = req.Notes
	existing.Amount = req.Amount

	if err := tx.Save(existing).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal update pembayaran")
	}

	saleMap := make(map[uint]bool)
	for _, s := range oldSaleIDs {
		saleMap[s] = true
	}
	for _, it := range paymentItems {
		saleMap[it.SaleID] = true
	}

	for sID := range saleMap {
		var totalPaid float64
		tx.Model(&model.ReceivablePaymentItem{}).
			Where("sale_id = ?", sID).
			Select("COALESCE(SUM(amount_paid), 0)").
			Scan(&totalPaid)

		var sale model.Sale
		tx.First(&sale, sID)

		status := "piutang"
		if totalPaid >= sale.GrandTotal {
			status = "lunas"
		}

		tx.Model(&sale).Updates(map[string]interface{}{
			"amount_paid":    totalPaid,
			"payment_status": status,
		})
	}

	tx.Commit()

	return repository.GetReceivablePaymentByID(id)
}

func DeleteReceivablePayment(id uint) error {
	existing, err := repository.GetReceivablePaymentByID(id)
	if err != nil {
		return errors.New("pembayaran tidak ditemukan")
	}

	tx := config.DB.Begin()
	var saleIDs []uint
	for _, it := range existing.Items {
		saleIDs = append(saleIDs, it.SaleID)
	}

	if err := tx.Where("receivable_payment_id = ?", id).Delete(&model.ReceivablePaymentItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&model.ReceivablePayment{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, sID := range saleIDs {
		var totalPaid float64
		tx.Model(&model.ReceivablePaymentItem{}).
			Where("sale_id = ?", sID).
			Select("COALESCE(SUM(amount_paid), 0)").
			Scan(&totalPaid)

		var sale model.Sale
		tx.First(&sale, sID)

		status := "piutang"
		if totalPaid >= sale.GrandTotal {
			status = "lunas"
		}

		tx.Model(&sale).Updates(map[string]interface{}{
			"amount_paid":    totalPaid,
			"payment_status": status,
		})
	}
	tx.Commit()
	return nil
}

func GetPaymentsBySaleID(saleID uint) ([]model.ReceivablePaymentItem, error) {
	return repository.GetPaymentsBySaleID(saleID)
}

func GetUnpaidSales(customerName string) ([]model.Sale, error) {
	return repository.GetUnpaidSales(customerName)
}

func GetReceivablePaymentByID(id uint) (*model.ReceivablePayment, error) {
	return repository.GetReceivablePaymentByID(id)
}

func FixMigrate() error {
	return repository.FixMigrate()
}

func SyncSales() error {
	return repository.SyncSales()
}
