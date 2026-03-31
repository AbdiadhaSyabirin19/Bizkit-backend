package service

import (
	"errors"
	"fmt"
	"strings"
	"time"
	_ "time/tzdata"

	"bizkit-backend/internal/model"
	"bizkit-backend/internal/repository"
)

type PromoItemRequest struct {
	RefType string `json:"ref_type"`
	RefID   uint   `json:"ref_id"`
	RefName string `json:"ref_name"`
}

type PromoSpecialPriceRequest struct {
	ProductID uint    `json:"product_id"`
	BuyPrice  float64 `json:"buy_price"`
}

type PromoRequest struct {
	Name          string                     `json:"name"`
	PromoType     string                     `json:"promo_type"`
	AppliesTo     string                     `json:"applies_to"`
	Condition     string                     `json:"condition"`
	MinQty        int                        `json:"min_qty"`
	MinTotal      float64                    `json:"min_total"`
	DiscountPct   float64                    `json:"discount_pct"`
	MaxDiscount   float64                    `json:"max_discount"`
	CutPrice      float64                    `json:"cut_price"`
	ActiveDays    string                     `json:"active_days"`
	StartTime     string                     `json:"start_time"`
	EndTime       string                     `json:"end_time"`
	StartDate     string                     `json:"start_date"`
	EndDate       string                     `json:"end_date"`
	VoucherType   string                     `json:"voucher_type"`
	VoucherCode   string                     `json:"voucher_code"`
	MaxUsage      int                        `json:"max_usage"`
	Status        string                     `json:"status"`
	Items         []PromoItemRequest         `json:"items"`
	SpecialPrices []PromoSpecialPriceRequest `json:"special_prices"`
}

func GetAllPromos(search string) ([]model.Promo, error) {
	return repository.GetAllPromos(search)
}

func GetPromoByID(id uint) (*model.Promo, error) {
	promo, err := repository.GetPromoByID(id)
	if err != nil {
		return nil, errors.New("Promo tidak ditemukan")
	}
	return promo, nil
}

func CreatePromo(req PromoRequest) (*model.Promo, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	startDate, _ := time.ParseInLocation("2006-01-02", req.StartDate, loc)
	endDate, _ := time.ParseInLocation("2006-01-02", req.EndDate, loc)

	if req.Status == "" {
		req.Status = "active"
	}

	promo := model.Promo{
		Name:        req.Name,
		PromoType:   req.PromoType,
		AppliesTo:   req.AppliesTo,
		Condition:   req.Condition,
		MinQty:      req.MinQty,
		MinTotal:    req.MinTotal,
		DiscountPct: req.DiscountPct,
		MaxDiscount: req.MaxDiscount,
		CutPrice:    req.CutPrice,
		ActiveDays:  req.ActiveDays,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		StartDate:   startDate,
		EndDate:     endDate,
		VoucherType: req.VoucherType,
		VoucherCode: req.VoucherCode,
		MaxUsage:    req.MaxUsage,
		Status:      req.Status,
	}

	var items []model.PromoItem
	for _, it := range req.Items {
		items = append(items, model.PromoItem{
			RefType: it.RefType,
			RefID:   it.RefID,
			RefName: it.RefName,
		})
	}

	var specialPrices []model.PromoSpecialPrice
	for _, sp := range req.SpecialPrices {
		specialPrices = append(specialPrices, model.PromoSpecialPrice{
			ProductID: sp.ProductID,
			BuyPrice:  sp.BuyPrice,
		})
	}

	// Generate vouchers
	var vouchers []model.PromoVoucher
	if req.VoucherType == "custom" && req.VoucherCode != "" {
		vouchers = append(vouchers, model.PromoVoucher{Code: req.VoucherCode})
	}

	err := repository.CreatePromo(&promo, items, specialPrices, vouchers)
	if err != nil {
		return nil, err
	}

	// Generate vouchers jika tipe generate
	if req.VoucherType == "generate" && req.MaxUsage > 0 {
		repository.GenerateVoucherCodes(promo.ID, req.MaxUsage)
	}

	return repository.GetPromoByID(promo.ID)
}

func UpdatePromo(id uint, req PromoRequest) (*model.Promo, error) {
	promo, err := repository.GetPromoByID(id)
	if err != nil {
		return nil, errors.New("Promo tidak ditemukan")
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	startDate, _ := time.ParseInLocation("2006-01-02", req.StartDate, loc)
	endDate, _ := time.ParseInLocation("2006-01-02", req.EndDate, loc)

	promo.Name = req.Name
	promo.PromoType = req.PromoType
	promo.AppliesTo = req.AppliesTo
	promo.Condition = req.Condition
	promo.MinQty = req.MinQty
	promo.MinTotal = req.MinTotal
	promo.DiscountPct = req.DiscountPct
	promo.MaxDiscount = req.MaxDiscount
	promo.CutPrice = req.CutPrice
	promo.ActiveDays = req.ActiveDays
	promo.StartTime = req.StartTime
	promo.EndTime = req.EndTime
	promo.StartDate = startDate
	promo.EndDate = endDate
	promo.VoucherType = req.VoucherType
	promo.VoucherCode = req.VoucherCode
	promo.MaxUsage = req.MaxUsage
	if req.Status != "" {
		promo.Status = req.Status
	}

	var items []model.PromoItem
	for _, it := range req.Items {
		items = append(items, model.PromoItem{
			RefType: it.RefType,
			RefID:   it.RefID,
			RefName: it.RefName,
		})
	}

	var specialPrices []model.PromoSpecialPrice
	for _, sp := range req.SpecialPrices {
		specialPrices = append(specialPrices, model.PromoSpecialPrice{
			ProductID: sp.ProductID,
			BuyPrice:  sp.BuyPrice,
		})
	}

	err = repository.UpdatePromo(promo, items, specialPrices)
	if err != nil {
		return nil, err
	}
	return repository.GetPromoByID(promo.ID)
}

func DeletePromo(id uint) error {
	_, err := repository.GetPromoByID(id)
	if err != nil {
		return errors.New("Promo tidak ditemukan")
	}
	return repository.DeletePromo(id)
}

func GetPromosByProductID(productID uint, categoryID *uint, brandID *uint) ([]model.Promo, error) {
	return repository.GetPromosByProductID(productID, categoryID, brandID)
}

// ── Check Request Structs ─────────────────────────────────────────────────

type CheckPromoItem struct {
	ProductID  uint    `json:"product_id"`
	CategoryID *uint   `json:"category_id"`
	BrandID    *uint   `json:"brand_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
}

type CheckPromoRequest struct {
	Items    []CheckPromoItem `json:"items"`
	Subtotal float64          `json:"subtotal"`
}

type PromoResult struct {
	PromoID        uint    `json:"promo_id"`
	Name           string  `json:"name"`
	PromoType      string  `json:"promo_type"`
	DiscountAmount float64 `json:"discount_amount"`
	Description    string  `json:"description"`
}

// ── CheckAutoPromos ───────────────────────────────────────────────────────
// Cek semua promo aktif yang berlaku untuk keranjang ini (tanpa voucher)

func CheckAutoPromos(req CheckPromoRequest) ([]PromoResult, error) {
	allPromos, err := repository.GetActivePromos()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var results []PromoResult

	for _, promo := range allPromos {
		// Skip promo berbasis voucher
		if promo.VoucherType != "none" && promo.VoucherType != "" {
			continue
		}

		if !isPromoValid(promo, now) {
			continue
		}

		if !isPromoApplicable(promo, req) {
			continue
		}

		discount := calcDiscount(promo, req)
		if discount <= 0 {
			continue
		}

		results = append(results, PromoResult{
			PromoID:        promo.ID,
			Name:           promo.Name,
			PromoType:      promo.PromoType,
			DiscountAmount: discount,
			Description:    buildPromoDescription(promo),
		})
	}

	return results, nil
}

// ── CheckVoucher ──────────────────────────────────────────────────────────
// Validasi kode voucher manual dari kasir

type CheckVoucherRequest struct {
	Code     string           `json:"code" binding:"required"`
	Items    []CheckPromoItem `json:"items"`
	Subtotal float64          `json:"subtotal"`
}

func CheckVoucher(req CheckVoucherRequest) (*PromoResult, error) {
	promo, voucher, err := repository.GetPromoByVoucherCode(req.Code)
	if err != nil {
		return nil, errors.New("Kode voucher tidak ditemukan")
	}

	if promo.Status != "active" {
		return nil, errors.New("Promo tidak aktif")
	}

	now := time.Now()
	if !isPromoValid(*promo, now) {
		return nil, errors.New("Promo sudah tidak berlaku")
	}

	if promo.MaxUsage > 0 && promo.UsedCount >= promo.MaxUsage {
		return nil, errors.New("Voucher sudah mencapai batas penggunaan")
	}

	// Untuk voucher per-code (generate), cek is_used
	if voucher != nil && voucher.IsUsed {
		return nil, errors.New("Kode voucher ini sudah digunakan")
	}

	checkReq := CheckPromoRequest{Items: req.Items, Subtotal: req.Subtotal}
	discount := calcDiscount(*promo, checkReq)
	if discount <= 0 {
		return nil, errors.New("Voucher tidak berlaku untuk transaksi ini")
	}

	return &PromoResult{
		PromoID:        promo.ID,
		Name:           promo.Name,
		PromoType:      promo.PromoType,
		DiscountAmount: discount,
		Description:    buildPromoDescription(*promo),
	}, nil
}

// ── Helper Functions ──────────────────────────────────────────────────────

func isPromoValid(promo model.Promo, now time.Time) bool {
	// 1. Pastikan waktu saat ini menggunakan zona Asia/Jakarta
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now = now.In(loc)

	// 2. Tentukan batasan waktu absolut (Tanggal + Jam)
	startTime := promo.StartTime
	if startTime == "" {
		startTime = "00:00"
	}
	endTime := promo.EndTime
	if endTime == "" {
		endTime = "23:59"
	}

	startStr := fmt.Sprintf("%s %s", promo.StartDate.Format("2006-01-02"), startTime)
	endStr := fmt.Sprintf("%s %s", promo.EndDate.Format("2006-01-02"), endTime)

	startDT, _ := time.ParseInLocation("2006-01-02 15:04", startStr, loc)
	endDT, _ := time.ParseInLocation("2006-01-02 15:04", endStr, loc)

	// Jika StartTime dan EndTime sama, asumsikan aktif 24 jam di hari tersebut
	if startTime == endTime && startTime != "00:00" {
		endDT = endDT.Add(24 * time.Hour)
	}

	if now.Before(startDT) || now.After(endDT) {
		return false
	}

	// 3. Cek hari aktif (1=Senin ... 7=Minggu)
	if promo.ActiveDays != "" {
		dayNum := int(now.Weekday())
		if dayNum == 0 {
			dayNum = 7 // Sunday = 7
		}
		dayStr := fmt.Sprintf("%d", dayNum)
		found := false
		for _, d := range strings.Split(promo.ActiveDays, ",") {
			if strings.TrimSpace(d) == dayStr {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func isPromoApplicable(promo model.Promo, req CheckPromoRequest) bool {
	// Cek kondisi minimum
	totalQty := 0
	for _, item := range req.Items {
		totalQty += item.Quantity
	}

	switch promo.Condition {
	case "qty":
		if totalQty < promo.MinQty {
			return false
		}
	case "total":
		if req.Subtotal < promo.MinTotal {
			return false
		}
	case "qty_and_total":
		if totalQty < promo.MinQty || req.Subtotal < promo.MinTotal {
			return false
		}
	case "qty_or_total":
		if totalQty < promo.MinQty && req.Subtotal < promo.MinTotal {
			return false
		}
	}

	// Cek applies_to
	if promo.AppliesTo == "all" {
		return true
	}

	// Cek apakah ada item yang match dengan promo items
	for _, item := range req.Items {
		for _, pi := range promo.Items {
			switch pi.RefType {
			case "product":
				if pi.RefID == item.ProductID {
					return true
				}
			case "category":
				if item.CategoryID != nil && pi.RefID == *item.CategoryID {
					return true
				}
			case "brand":
				if item.BrandID != nil && pi.RefID == *item.BrandID {
					return true
				}
			}
		}
	}

	return false
}

func calcDiscount(promo model.Promo, req CheckPromoRequest) float64 {
	// Hitung subtotal hanya untuk item yang masuk dalam cakupan promo
	applicableSubtotal := 0.0
	if promo.AppliesTo == "all" {
		applicableSubtotal = req.Subtotal
	} else {
		for _, item := range req.Items {
			isMatch := false
			for _, pi := range promo.Items {
				if (pi.RefType == "product" && pi.RefID == item.ProductID) ||
					(pi.RefType == "category" && item.CategoryID != nil && pi.RefID == *item.CategoryID) ||
					(pi.RefType == "brand" && item.BrandID != nil && pi.RefID == *item.BrandID) {
					isMatch = true
					break
				}
			}
			if isMatch {
				applicableSubtotal += item.Price * float64(item.Quantity)
			}
		}
	}

	if applicableSubtotal <= 0 {
		return 0
	}

	switch promo.PromoType {
	case "discount":
		discount := applicableSubtotal * (promo.DiscountPct / 100)
		if promo.MaxDiscount > 0 && discount > promo.MaxDiscount {
			discount = promo.MaxDiscount
		}
		return discount
	case "cut_price":
		// Potongan harga tetap (cut_price) biasanya langsung memotong total akhir,
		// tapi kita batasi jangan sampai melebihi subtotal item yang promo
		discount := promo.CutPrice
		if discount > applicableSubtotal {
			discount = applicableSubtotal
		}
		return discount
	case "special_price":
		// Hitung selisih harga normal vs special price per item
		var saving float64
		for _, item := range req.Items {
			for _, sp := range promo.SpecialPrices {
				if sp.ProductID == item.ProductID {
					diff := item.Price - sp.BuyPrice
					if diff > 0 {
						saving += diff * float64(item.Quantity)
					}
				}
			}
		}
		return saving
	}
	return 0
}

func buildPromoDescription(promo model.Promo) string {
	switch promo.PromoType {
	case "discount":
		desc := fmt.Sprintf("Diskon %.0f%%", promo.DiscountPct)
		if promo.MaxDiscount > 0 {
			desc += fmt.Sprintf(" (maks Rp %.0f)", promo.MaxDiscount)
		}
		return desc
	case "cut_price":
		return fmt.Sprintf("Potongan Rp %.0f", promo.CutPrice)
	case "special_price":
		return "Harga spesial"
	}
	return promo.Name
}
