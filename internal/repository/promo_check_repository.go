package repository

import (
	"bizkit-backend/config"
	"bizkit-backend/internal/model"
)

// GetActivePromos — ambil semua promo status active beserta relasi
func GetActivePromos() ([]model.Promo, error) {
	var promos []model.Promo
	err := config.DB.
		Preload("Items").
		Preload("SpecialPrices").
		Where("status = ?", "active").
		Find(&promos).Error
	return promos, err
}

// GetPromoByVoucherCode — cari promo berdasarkan kode voucher
// Mengembalikan (promo, voucher, error)
// - Untuk VoucherType "custom": voucher = nil (kode ada di promo.VoucherCode)
// - Untuk VoucherType "generate": voucher = record PromoVoucher spesifik
func GetPromoByVoucherCode(code string) (*model.Promo, *model.PromoVoucher, error) {
	// Cek voucher type "custom" — kode langsung di field promo
	var promoCustom model.Promo
	err := config.DB.
		Preload("Items").
		Preload("SpecialPrices").
		Where("voucher_code = ? AND voucher_type = ? AND status = ?", code, "custom", "active").
		First(&promoCustom).Error
	if err == nil {
		return &promoCustom, nil, nil
	}

	// Cek voucher type "generate" — kode ada di tabel promo_vouchers
	var voucher model.PromoVoucher
	err = config.DB.Where("code = ?", code).First(&voucher).Error
	if err != nil {
		return nil, nil, err
	}

	var promo model.Promo
	err = config.DB.
		Preload("Items").
		Preload("SpecialPrices").
		Where("id = ? AND status = ?", voucher.PromoID, "active").
		First(&promo).Error
	if err != nil {
		return nil, nil, err
	}

	return &promo, &voucher, nil
}
