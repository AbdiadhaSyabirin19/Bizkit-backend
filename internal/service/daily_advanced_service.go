package service

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"bizkit-backend/internal/model"
	"bizkit-backend/internal/repository"
)

type ProductStat struct {
	Name      string  `json:"name"`
	Qty       int     `json:"qty"`
	NotaCount int     `json:"nota_count"`
	Omzet     float64 `json:"omzet"`
}

type CategoryStat struct {
	Name      string  `json:"name"`
	Qty       int     `json:"qty"`
	NotaCount int     `json:"nota_count"`
	Omzet     float64 `json:"omzet"`
}

type ShiftCashFlow struct {
	ShiftID       uint       `json:"shift_id"`
	UserName      string     `json:"user_name"`
	StartTime     time.Time  `json:"start_time"`
	EndTime       *time.Time `json:"end_time"`
	CashSales     float64    `json:"cash_sales"`
	TotalKasMasuk float64    `json:"total_kas_masuk"`
	TotalKasTunai float64    `json:"total_kas_tunai"`
}

type PaymentItem struct {
	Name  string  `json:"name"`
	Total float64 `json:"total"`
}

// Rekapitulasi per grup (hari/bulan/minggu)
type GroupedSale struct {
	Label string  `json:"label"` // e.g. "01 Mar", "Januari", "Minggu 1"
	Date  string  `json:"date"`  // sortable date e.g. "2026-03-01"
	N     int     `json:"n"`     // Nota count
	Q     int     `json:"q"`     // Qty sum
	Omzet float64 `json:"omzet"` // Omzet sum
}

type ChartDataPoint struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

func getWeekNumber(date time.Time) int {
	_, week := date.ISOWeek()
	return week
}

func GetDailyAdvancedReport(startStr, endStr, mode, yearlyBreakdown string) (map[string]interface{}, error) {
	if startStr == "" || endStr == "" {
		return nil, errors.New("start_date dan end_date wajib diisi")
	}

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		return nil, errors.New("Format start_date tidak valid")
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		return nil, errors.New("Format end_date tidak valid")
	}

	endFull := end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// 1. Ambil semua sales dalam periode dengan relasi lengkap
	sales, err := repository.GetSalesByPeriod(start, endFull)
	if err != nil {
		return nil, err
	}

	// 2. Jika mode daily, ambil shifts. Jika bukan, shifts kosongkan saja untuk menghemat processing
	var shifts []model.Shift
	if mode == "daily" {
		shifts, _ = repository.GetShiftsByPeriod(start, endFull)
	}

	// Hitung Aggregates Stats (Sama untuk semua mode)
	var totalOmzet float64
	var totalQty int
	paymentSummary := map[string]float64{}
	productMap := map[string]*ProductStat{}
	productNota := map[string]map[uint]bool{}
	categoryMap := map[string]*CategoryStat{}
	categoryNota := map[string]map[uint]bool{}

	// Untuk data charts (hanya non-daily)
	hourlyQty := map[int]struct {
		TotalQty  int
		DaysCount map[string]bool
	}{}
	dayOfWeekQty := map[time.Weekday]struct {
		TotalQty   int
		WeeksCount map[string]bool
	}{}

	for i := 0; i < 24; i++ {
		hourlyQty[i] = struct {
			TotalQty  int
			DaysCount map[string]bool
		}{DaysCount: map[string]bool{}}
	}
	for i := time.Sunday; i <= time.Saturday; i++ {
		dayOfWeekQty[i] = struct {
			TotalQty   int
			WeeksCount map[string]bool
		}{WeeksCount: map[string]bool{}}
	}

	// Agregasi Grouping (untuk non-daily)
	groupMap := map[string]*GroupedSale{}

	for _, sale := range sales {
		saleDate := sale.CreatedAt
		totalOmzet += sale.GrandTotal

		methodName := sale.PaymentMethod.Name
		if methodName == "" {
			methodName = "Lainnya"
		}
		paymentSummary[methodName] += sale.GrandTotal

		saleQty := 0
		for _, item := range sale.Items {
			saleQty += item.Quantity
			pName := item.Product.Name
			if pName != "" {
				if _, ok := productMap[pName]; !ok {
					productMap[pName] = &ProductStat{Name: pName}
					productNota[pName] = map[uint]bool{}
				}
				productMap[pName].Qty += item.Quantity
				productMap[pName].Omzet += item.Subtotal
				productNota[pName][sale.ID] = true
			}

			catName := "Tanpa Kategori"
			if item.Product.Category != nil && item.Product.Category.Name != "" {
				catName = item.Product.Category.Name
			}
			if _, ok := categoryMap[catName]; !ok {
				categoryMap[catName] = &CategoryStat{Name: catName}
				categoryNota[catName] = map[uint]bool{}
			}
			categoryMap[catName].Qty += item.Quantity
			categoryMap[catName].Omzet += item.Subtotal
			categoryNota[catName][sale.ID] = true
		}
		totalQty += saleQty

		// Perekaman jam dan hari untuk chart (hanya relevan untuk mode non-daily)
		if mode != "daily" {
			hour := saleDate.Hour()
			dateKey := saleDate.Format("2006-01-02")
			hData := hourlyQty[hour]
			hData.TotalQty += saleQty
			hData.DaysCount[dateKey] = true
			hourlyQty[hour] = hData

			weekday := saleDate.Weekday()
			year, week := saleDate.ISOWeek()
			weekKey := fmt.Sprintf("%d-%d", year, week)
			dData := dayOfWeekQty[weekday]
			dData.TotalQty += saleQty
			dData.WeeksCount[weekKey] = true
			dayOfWeekQty[weekday] = dData

			// Logic Grouping
			var groupKey, groupLabel string
			if mode == "weekly" || mode == "monthly" {
				groupKey = saleDate.Format("2006-01-02") // Grup by hari
				if mode == "weekly" {
					daysIndo := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
					groupLabel = daysIndo[saleDate.Weekday()]
				} else {
					groupLabel = fmt.Sprintf("%02d/%02d", saleDate.Day(), saleDate.Month())
				}
			} else if mode == "yearly" {
				if yearlyBreakdown == "weekly" {
					year, week := saleDate.ISOWeek()
					groupKey = fmt.Sprintf("%04d-W%02d", year, week)
					groupLabel = fmt.Sprintf("Minggu %d", week)
				} else {
					// default monthly breakdown for yearly
					groupKey = saleDate.Format("2006-01") // Grup by bulan
					monthNames := []string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
					groupLabel = monthNames[saleDate.Month()]
				}
			}

			if _, ok := groupMap[groupKey]; !ok {
				groupMap[groupKey] = &GroupedSale{
					Label: groupLabel,
					Date:  groupKey,
				}
			}
			groupMap[groupKey].N++
			groupMap[groupKey].Q += saleQty
			groupMap[groupKey].Omzet += sale.GrandTotal
		}
	}

	// Hitung Produk & Kategori
	var productStats []ProductStat
	for name, stat := range productMap {
		stat.NotaCount = len(productNota[name])
		productStats = append(productStats, *stat)
	}

	var categoryStats []CategoryStat
	for name, stat := range categoryMap {
		stat.NotaCount = len(categoryNota[name])
		categoryStats = append(categoryStats, *stat)
	}

	var paymentItems []PaymentItem
	for name, total := range paymentSummary {
		paymentItems = append(paymentItems, PaymentItem{Name: name, Total: total})
	}

	result := map[string]interface{}{
		"period": map[string]string{
			"start": startStr,
			"end":   endStr,
		},
		"total_transaksi": len(sales), // N total
		"total_qty":       totalQty,   // Q total
		"total_omzet":     totalOmzet,
		"payment_summary": paymentItems,
		"product_stats":   productStats,
		"category_stats":  categoryStats,
	}

	if mode == "daily" {
		// Shift Cash Flow cuma buat daily
		var shiftCashFlows []ShiftCashFlow
		for _, shift := range shifts {
			scf := ShiftCashFlow{
				ShiftID:   shift.ID,
				UserName:  shift.User.Name,
				StartTime: shift.StartTime,
				EndTime:   shift.EndTime,
			}
			var shiftEnd time.Time
			if shift.EndTime != nil {
				shiftEnd = *shift.EndTime
			} else {
				shiftEnd = endFull
			}
			for _, sale := range sales {
				saleTime := sale.CreatedAt
				if saleTime.After(shift.StartTime) && (saleTime.Before(shiftEnd) || saleTime.Equal(shiftEnd)) &&
					sale.UserID == shift.UserID {
					methodName := sale.PaymentMethod.Name
					if methodName == "Tunai" || methodName == "Cash" || methodName == "tunai" || methodName == "cash" {
						scf.CashSales += sale.GrandTotal
					}
					scf.TotalKasMasuk += sale.GrandTotal
				}
			}
			scf.TotalKasTunai = shift.CashIn + scf.CashSales
			shiftCashFlows = append(shiftCashFlows, scf)
		}

		result["sales"] = sales
		result["shift_cash_flows"] = shiftCashFlows

	} else {
		// Kumpulkan dan sort groups
		var groupedSales []GroupedSale
		for _, g := range groupMap {
			groupedSales = append(groupedSales, *g)
		}
		sort.Slice(groupedSales, func(i, j int) bool {
			return groupedSales[i].Date < groupedSales[j].Date
		})

		// Summary Rata-rata
		groupCount := len(groupedSales)
		if groupCount == 0 {
			groupCount = 1 // cegah div by zero
		}
		totalN := float64(len(sales))
		totalQ := float64(totalQty)
		if totalN == 0 {
			totalN = 1
		}
		if totalQ == 0 {
			totalQ = 1
		}

		summary := map[string]float64{
			"avg_nota":         float64(len(sales)) / float64(groupCount),
			"avg_qty":          float64(totalQty) / float64(groupCount),
			"avg_omzet":        totalOmzet / float64(groupCount),
			"avg_qty_per_nota": float64(totalQty) / totalN,
			"avg_omz_per_nota": totalOmzet / totalN,
			"avg_omz_per_qty":  totalOmzet / totalQ,
		}

		// Chart 1: Trend Qty
		var chartTrend []ChartDataPoint
		for _, g := range groupedSales {
			chartTrend = append(chartTrend, ChartDataPoint{Label: g.Label, Value: float64(g.Q)})
		}

		// Chart 2: Analisis Perhari
		var chartDay []ChartDataPoint
		days := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
		for i := 0; i < 7; i++ {
			idx := time.Weekday(i)
			dData := dayOfWeekQty[idx]
			count := len(dData.WeeksCount)
			if count == 0 {
				count = 1
			}
			chartDay = append(chartDay, ChartDataPoint{
				Label: days[i],
				Value: float64(dData.TotalQty) / float64(count),
			})
		}
		// Referensi pake Senin - Minggu. Jadi kita rotate
		// chartDay => Minggu di awal. Kita pindah Minggu ke akhir
		if len(chartDay) == 7 {
			chartDay = append(chartDay[1:], chartDay[0])
		}

		// Chart 3: Analisis Perjam
		var chartHour []ChartDataPoint
		for i := 0; i < 24; i++ {
			hData := hourlyQty[i]
			count := len(hData.DaysCount)
			if count == 0 {
				count = 1
			}
			chartHour = append(chartHour, ChartDataPoint{
				Label: fmt.Sprintf("%02d:00", i),
				Value: float64(hData.TotalQty) / float64(count),
			})
		}

		result["grouped_sales"] = groupedSales
		result["summary"] = summary
		result["chart_trend"] = chartTrend
		result["chart_day"] = chartDay
		result["chart_hour"] = chartHour
	}

	return result, nil
}
