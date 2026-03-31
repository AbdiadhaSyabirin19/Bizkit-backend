package service

import (
	"errors"
	"time"

	"bizkit-backend/internal/model"
	"bizkit-backend/internal/repository"
)

func parsePeriod(startStr, endStr string) (time.Time, time.Time, error) {
	if startStr == "" || endStr == "" {
		return time.Time{}, time.Time{}, errors.New("start_date dan end_date wajib diisi")
	}

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("Format start_date tidak valid")
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("Format end_date tidak valid")
	}

	// Set end ke akhir hari
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	return start, end, nil
}

func GetSalesReport(startStr, endStr string, onlyDiscounted bool) (map[string]interface{}, error) {
	start, end, err := parsePeriod(startStr, endStr)
	if err != nil {
		return nil, err
	}

	sales, err := repository.GetSalesByPeriod(start, end)
	if err != nil {
		return nil, err
	}

	// Filter jika hanya minta data diskon
	var filteredSales []model.Sale
	if onlyDiscounted {
		for _, sale := range sales {
			if sale.DiscountTotal > 0 || sale.ManualDiscount > 0 {
				filteredSales = append(filteredSales, sale)
			}
		}
	} else {
		filteredSales = sales
	}

	// Rekap total
	var totalOmzet, totalDiskon float64
	paymentSummary := map[string]float64{}

	for _, sale := range filteredSales {
		totalOmzet += sale.GrandTotal
		totalDiskon += sale.DiscountTotal + sale.ManualDiscount
		paymentSummary[sale.PaymentMethod.Name] += sale.GrandTotal
	}

	return map[string]interface{}{
		"period": map[string]string{
			"start": startStr,
			"end":   endStr,
		},
		"total_transaksi": len(filteredSales),
		"total_omzet":     totalOmzet,
		"total_diskon":    totalDiskon,
		"payment_summary": paymentSummary,
		"sales":           filteredSales,
	}, nil
}

func GetTrendReport(startStr, endStr string) (map[string]interface{}, error) {
	start, end, err := parsePeriod(startStr, endStr)
	if err != nil {
		return nil, err
	}

	items, err := repository.GetSaleItemsByPeriod(start, end)
	if err != nil {
		return nil, err
	}

	// Rekap per produk
	productStats := map[string]map[string]interface{}{}
	categoryStats := map[string]map[string]interface{}{}

	for _, item := range items {
		// Per produk
		productName := item.Product.Name
		if productName == "" {
			continue
		}
		if _, ok := productStats[productName]; !ok {
			productStats[productName] = map[string]interface{}{
				"name":  productName,
				"qty":   0,
				"omzet": 0.0,
			}
		}
		productStats[productName]["qty"] = productStats[productName]["qty"].(int) + item.Quantity
		productStats[productName]["omzet"] = productStats[productName]["omzet"].(float64) + item.Subtotal

		// Per kategori — cek nil dulu!
		categoryName := ""
		if item.Product.Category.Name != "" {
			categoryName = item.Product.Category.Name
		} else {
			categoryName = "Tanpa Kategori"
		}

		if _, ok := categoryStats[categoryName]; !ok {
			categoryStats[categoryName] = map[string]interface{}{
				"name":  categoryName,
				"qty":   0,
				"omzet": 0.0,
			}
		}
		categoryStats[categoryName]["qty"] = categoryStats[categoryName]["qty"].(int) + item.Quantity
		categoryStats[categoryName]["omzet"] = categoryStats[categoryName]["omzet"].(float64) + item.Subtotal
	}

	// Convert map ke slice
	var productList []interface{}
	for _, v := range productStats {
		productList = append(productList, v)
	}

	var categoryList []interface{}
	for _, v := range categoryStats {
		categoryList = append(categoryList, v)
	}

	return map[string]interface{}{
		"period": map[string]string{
			"start": startStr,
			"end":   endStr,
		},
		"product_stats":  productList,
		"category_stats": categoryList,
	}, nil
}

func GetAttendanceReport(dateStr string) (map[string]interface{}, error) {
	var date time.Time
	var err error

	if dateStr == "" {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, errors.New("Format tanggal tidak valid")
		}
	}

	attendances, err := repository.GetAttendanceByDate(date)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"date":        date.Format("2006-01-02"),
		"total":       len(attendances),
		"attendances": attendances,
	}, nil
}

func GetShiftReport(startStr, endStr string) (map[string]interface{}, error) {
	start, end, err := parsePeriod(startStr, endStr)
	if err != nil {
		return nil, err
	}

	shifts, err := repository.GetShiftsByPeriod(start, end)
	if err != nil {
		return nil, err
	}

	// Hitung total shift
	totalShift := len(shifts)

	// Total kas masuk = jumlah semua CashIn (modal buka shift)
	var totalCashIn float64
	for _, shift := range shifts {
		totalCashIn += shift.CashIn
	}

	// Total kas keluar = jumlah semua CashOut (uang yang disetor/keluar saat tutup shift)
	var totalCashOut float64
	for _, shift := range shifts {
		totalCashOut += shift.CashOut
	}

	// Total penjualan dari semua shift (ambil dari sales dalam rentang periode)
	sales, err := repository.GetSalesByPeriod(start, end)
	if err != nil {
		return nil, err
	}
	var totalSales float64
	for _, s := range sales {
		totalSales += s.GrandTotal
	}

	// Saldo akhir = total kas masuk + total penjualan - total kas keluar
	saldoAkhir := totalCashIn + totalSales - totalCashOut

	return map[string]interface{}{
		"period": map[string]string{
			"start": startStr,
			"end":   endStr,
		},
		"total_shift":    totalShift,
		"total_cash_in":  totalCashIn,
		"total_cash_out": totalCashOut,
		"total_sales":    totalSales,
		"saldo_akhir":    saldoAkhir,
		"shifts":         shifts,
	}, nil
}

func GetDashboardSummary() (map[string]interface{}, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// 1. Sales Today
	sales, _ := repository.GetSalesByPeriod(todayStart, todayEnd)
	var omzetToday float64
	for _, s := range sales {
		omzetToday += s.GrandTotal
	}

	// 2. Transaksi Today
	transCount := len(sales)

	// 3. Active Shift
	activeShift, _ := repository.GetActiveShiftAny()

	// 4. Low Stock Products (misal stock <= 5)
	products, _ := repository.GetAllProducts("")
	lowStockCount := 0
	for _, p := range products {
		if p.Stock <= 5 {
			lowStockCount++
		}
	}

	return map[string]interface{}{
		"date":             now.Format("2006-01-02"),
		"omzet_today":      omzetToday,
		"trans_count":      transCount,
		"has_active_shift": activeShift != nil && activeShift.ID != 0,
		"active_shift":     activeShift,
		"low_stock_count":  lowStockCount,
	}, nil
}
