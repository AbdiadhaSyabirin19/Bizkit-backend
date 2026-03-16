package repository

import (
	"bizkit-backend/config"
	"bizkit-backend/internal/model"
	"time"

	"gorm.io/gorm"
)

func GetAllSales(startDate, endDate string) ([]model.Sale, error) {
	var sales []model.Sale
	query := config.DB.Model(&model.Sale{}).
		Preload("User").
		Preload("PaymentMethod").
		Preload("Promo").
		Preload("Items.Product").
		Preload("Items.Variants.VariantOption")

	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	result := query.Order("created_at DESC").Find(&sales)
	return sales, result.Error
}

func GetSaleByID(id uint) (*model.Sale, error) {
	var sale model.Sale
	result := config.DB.
		Preload("User").
		Preload("PaymentMethod").
		Preload("Promo").
		Preload("Items.Product").
		Preload("Items.Variants.VariantOption").
		First(&sale, id)
	return &sale, result.Error
}

func GetSaleByInvoice(invoice string) (*model.Sale, error) {
	var sale model.Sale
	result := config.DB.Where("invoice_number = ?", invoice).First(&sale)
	return &sale, result.Error
}

func CreateSale(sale *model.Sale) error {
	return config.DB.Create(sale).Error
}

func UpdateSale(sale *model.Sale) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// Hapus items lama (soft delete)
		if err := tx.Where("sale_id = ?", sale.ID).Delete(&model.SaleItem{}).Error; err != nil {
			return err
		}
		// Hapus item variants lama (soft delete) - ini agak tricky karena gorm soft delete butuh ID
		// Tapi kita bisa hapus berdasarkan sale_item_id yang dimiliki sale ini
		if err := tx.Exec("UPDATE sale_item_variants SET deleted_at = ? WHERE sale_item_id IN (SELECT id FROM sale_items WHERE sale_id = ?)", time.Now(), sale.ID).Error; err != nil {
			return err
		}

		// Simpan perubahan header
		if err := tx.Save(sale).Error; err != nil {
			return err
		}

		return nil
	})
}

func DeleteSale(id uint) error {
	return config.DB.Delete(&model.Sale{}, id).Error
}

func GetDailySales(date time.Time, source string) ([]model.Sale, error) {
	var sales []model.Sale
	start := date.Format("2006-01-02") + " 00:00:00"
	end := date.Format("2006-01-02") + " 23:59:59"

	query := config.DB.
		Preload("User").
		Preload("PaymentMethod").
		Preload("Items.Product").
		Where("created_at BETWEEN ? AND ?", start, end)

	if source != "" {
		query = query.Where("source = ?", source)
	}

	result := query.Order("created_at DESC").Find(&sales)
	return sales, result.Error
}
func GetSalesByTimeRange(start, end time.Time, userID uint) ([]model.Sale, error) {
    var sales []model.Sale
    err := config.DB.
        Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, start, end).
        Find(&sales).Error
    return sales, err
}

// GetSalesByPeriod — semua sales dalam rentang waktu, dengan relasi lengkap
func GetSalesByPeriod(start, end time.Time) ([]model.Sale, error) {
	var sales []model.Sale
	result := config.DB.
		Preload("User").
		Preload("PaymentMethod").
		Preload("Promo").
		Preload("Items.Product.Category").
		Preload("Items.Variants.VariantOption").
		Where("created_at BETWEEN ? AND ?", start, end).
		Order("created_at DESC").
		Find(&sales)
	return sales, result.Error
}

// GetSaleItemsByPeriod — semua sale items dalam rentang waktu, untuk trend report
func GetSaleItemsByPeriod(start, end time.Time) ([]model.SaleItem, error) {
	var items []model.SaleItem
	result := config.DB.
		Preload("Product.Category").
		Joins("JOIN sales ON sale_items.sale_id = sales.id").
		Where("sales.created_at BETWEEN ? AND ? AND sales.deleted_at IS NULL", start, end).
		Find(&items)
	return items, result.Error
}