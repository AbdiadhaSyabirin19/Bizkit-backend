package service

import (
	"errors"
	"fmt"
	"time"

	"bizkit-backend/config"
	"bizkit-backend/internal/model"
	"bizkit-backend/internal/repository"
)

type SaleItemVariantRequest struct {
	VariantOptionID uint `json:"variant_option_id" binding:"required"`
}

type SaleItemRequest struct {
	ProductID uint                     `json:"product_id" binding:"required"`
	Quantity  int                      `json:"quantity" binding:"required"`
	Discount  float64                  `json:"discount"`
	Variants  []SaleItemVariantRequest `json:"variants"`
}

type SaleRequest struct {
	PaymentMethodID uint              `json:"payment_method_id" binding:"required"`
	PriceCategoryID *uint             `json:"price_category_id"`
	PromoID         *uint             `json:"promo_id"`
	CustomerName    string            `json:"customer_name" binding:"required"`
	Source          string            `json:"source"`
	ManualDiscount  float64           `json:"manual_discount"`
	AdditionalFee   float64           `json:"additional_fee"`
	Items           []SaleItemRequest `json:"items" binding:"required"`
}

func generateInvoiceNumber() string {
	now := time.Now()
	return fmt.Sprintf("INV-%s-%d", now.Format("20060102"), now.UnixNano()%10000)
}

func calculateSaleDetails(req SaleRequest) ([]model.SaleItem, float64, float64, float64, float64, float64, error) {
	var subtotal float64
	var saleItems []model.SaleItem

	for _, itemReq := range req.Items {
		product, err := repository.GetProductByID(itemReq.ProductID)
		if err != nil {
			return nil, 0, 0, 0, 0, 0, fmt.Errorf("Produk ID %d tidak ditemukan", itemReq.ProductID)
		}

		// Ambil harga (check multi-harga if category provided)
		basePrice := product.Price
		if req.PriceCategoryID != nil {
			var customPrice model.ProductPrice
			err := config.DB.Where("product_id = ? AND price_category_id = ?", product.ID, *req.PriceCategoryID).First(&customPrice).Error
			if err == nil && customPrice.Price > 0 {
				basePrice = customPrice.Price
			}
		}

		itemSubtotal := basePrice * float64(itemReq.Quantity)

		var saleItemVariants []model.SaleItemVariant
		for _, v := range itemReq.Variants {
			var option model.VariantOption
			if err := repository.GetVariantOptionByID(v.VariantOptionID, &option); err != nil {
				return nil, 0, 0, 0, 0, 0, fmt.Errorf("Variant option ID %d tidak ditemukan", v.VariantOptionID)
			}
			itemSubtotal += option.AdditionalPrice * float64(itemReq.Quantity)
			saleItemVariants = append(saleItemVariants, model.SaleItemVariant{
				VariantOptionID: v.VariantOptionID,
				AdditionalPrice: option.AdditionalPrice,
			})
		}

		itemSubtotal -= itemReq.Discount

		subtotal += itemSubtotal
		saleItems = append(saleItems, model.SaleItem{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			BasePrice: basePrice,
			Discount:  itemReq.Discount,
			Subtotal:  itemSubtotal,
			Variants:  saleItemVariants,
			Product:   *product, // Penting untuk promo check
		})
	}

	discountTotal, err := calculatePromoDiscount(req.PromoID, subtotal, saleItems)
	if err != nil {
		return nil, 0, 0, 0, 0, 0, err
	}

	grandTotal := (subtotal - discountTotal - req.ManualDiscount) + req.AdditionalFee
	return saleItems, subtotal, discountTotal, req.ManualDiscount, req.AdditionalFee, grandTotal, nil
}

func calculatePromoDiscount(promoID *uint, subtotal float64, items []model.SaleItem) (float64, error) {
	if promoID == nil {
		return 0, nil
	}

	promo, err := repository.GetPromoByID(*promoID)
	if err != nil {
		return 0, errors.New("Promo tidak ditemukan")
	}

	// Persiapkan request untuk pengecekan promo
	var checkItems []CheckPromoItem
	for _, it := range items {
		checkItems = append(checkItems, CheckPromoItem{
			ProductID:  it.ProductID,
			CategoryID: it.Product.CategoryID,
			BrandID:    it.Product.BrandID,
			Quantity:   it.Quantity,
			Price:      it.BasePrice + calculateVariantExtra(it),
		})
	}

	checkReq := CheckPromoRequest{
		Items:    checkItems,
		Subtotal: subtotal,
	}

	now := time.Now()
	if promo.Status != "active" {
		return 0, errors.New("Promo sudah tidak aktif")
	}
	if !isPromoValid(*promo, now) {
		return 0, errors.New("Promo sedang tidak berlaku (cek tanggal/jam/hari)")
	}
	if !isPromoApplicable(*promo, checkReq) {
		return 0, errors.New("Syarat promo tidak terpenuhi (min. qty/total atau produk tidak sesuai)")
	}

	if promo.MaxUsage > 0 && promo.UsedCount >= promo.MaxUsage {
		return 0, errors.New("Promo sudah mencapai batas penggunaan")
	}

	return calcDiscount(*promo, checkReq), nil
}

func calculateVariantExtra(item model.SaleItem) float64 {
	var extra float64
	for _, v := range item.Variants {
		extra += v.AdditionalPrice
	}
	return extra
}

func CreateSale(req SaleRequest, userID uint) (*model.Sale, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("Item transaksi tidak boleh kosong")
	}

	saleItems, subtotal, discountTotal, manualDiscount, additionalFee, grandTotal, err := calculateSaleDetails(req)
	if err != nil {
		return nil, err
	}

	sale := model.Sale{
		InvoiceNumber:   generateInvoiceNumber(),
		UserID:          userID,
		CustomerName:    req.CustomerName,
		PaymentMethodID: req.PaymentMethodID,
		PriceCategoryID: req.PriceCategoryID,
		PromoID:         req.PromoID,
		Subtotal:        subtotal,
		DiscountTotal:   discountTotal,
		ManualDiscount:  manualDiscount,
		AdditionalFee:   additionalFee,
		GrandTotal:      grandTotal,
		Source:          req.Source,
		Items:           saleItems,
	}

	if sale.Source == "" {
		sale.Source = "dashboard"
	}

	if err := repository.CreateSale(&sale); err != nil {
		return nil, errors.New("Gagal membuat transaksi")
	}

	if req.PromoID != nil {
		repository.UpdatePromoUsage(*req.PromoID)
	}

	result, _ := repository.GetSaleByID(sale.ID)
	return result, nil
}

func UpdateSale(id uint, req SaleRequest) (*model.Sale, error) {
	existing, err := repository.GetSaleByID(id)
	if err != nil {
		return nil, errors.New("Transaksi tidak ditemukan")
	}

	saleItems, subtotal, discountTotal, manualDiscount, additionalFee, grandTotal, err := calculateSaleDetails(req)
	if err != nil {
		return nil, err
	}

	existing.CustomerName = req.CustomerName
	existing.PaymentMethodID = req.PaymentMethodID
	existing.PriceCategoryID = req.PriceCategoryID
	existing.PromoID = req.PromoID
	existing.Subtotal = subtotal
	existing.DiscountTotal = discountTotal
	existing.ManualDiscount = manualDiscount
	existing.AdditionalFee = additionalFee
	existing.GrandTotal = grandTotal
	existing.Source = req.Source
	existing.Items = saleItems

	if err := repository.UpdateSale(existing); err != nil {
		return nil, errors.New("Gagal memperbarui transaksi: " + err.Error())
	}

	result, _ := repository.GetSaleByID(id)
	return result, nil
}

func DeleteSale(id uint) error {
	return repository.DeleteSale(id)
}

func GetAllSales(startDate, endDate string) ([]model.Sale, error) {
	return repository.GetAllSales(startDate, endDate)
}

func GetSaleByID(id uint) (*model.Sale, error) {
	sale, err := repository.GetSaleByID(id)
	if err != nil {
		return nil, errors.New("Transaksi tidak ditemukan")
	}
	return sale, nil
}

func GetDailySales(dateStr, source string) (map[string]interface{}, error) {
	var date time.Time
	var err error

	if dateStr == "" {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, errors.New("Format tanggal tidak valid, gunakan YYYY-MM-DD")
		}
	}

	sales, err := repository.GetDailySales(date, source)
	if err != nil {
		return nil, err
	}

	// Rekap metode pembayaran
	paymentSummary := map[string]float64{}
	var totalOmzet float64
	var totalQty int

	for _, sale := range sales {
		totalOmzet += sale.GrandTotal
		paymentSummary[sale.PaymentMethod.Name] += sale.GrandTotal
		for _, item := range sale.Items {
			totalQty += item.Quantity
		}
	}

	return map[string]interface{}{
		"date":            date.Format("2006-01-02"),
		"total_transaksi": len(sales),
		"total_qty":       totalQty,
		"total_omzet":     totalOmzet,
		"payment_summary": paymentSummary,
		"sales":           sales,
	}, nil
}