package main

import (
	"fmt"
	"time"

	"bizkit-backend/config"
	"bizkit-backend/internal/model"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	godotenv.Load()
	config.ConnectDB()

	fmt.Println("🌱 Mulai seeding data...")

	// ==================
	// ROLES
	// ==================
	roles := []model.Role{
		{Name: "Owner"},
		{Name: "Kasir"},
		{Name: "Admin"},
		{Name: "Supervisor"},
	}
	for i := range roles {
		config.DB.FirstOrCreate(&roles[i], model.Role{Name: roles[i].Name})
	}
	fmt.Println("✅ Roles selesai")

	// ==================
	// USERS
	// ==================
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	hashedStr := string(hashed)

	users := []struct {
		Name     string
		Username string
		RoleIdx  int
	}{
		{"Admin Bizkit", "admin", 0},
		{"Budi Kasir", "budi", 1},
		{"Siti Admin", "siti", 2},
		{"Andi Supervisor", "andi", 3},
	}
	for _, u := range users {
		var existing model.User
		if config.DB.Where("username = ?", u.Username).First(&existing).Error != nil {
			user := model.User{
				Name:     u.Name,
				Username: u.Username,
				Password: hashedStr,
				RoleID:   &roles[u.RoleIdx].ID,
			}
			config.DB.Create(&user)
		}
	}
	fmt.Println("✅ Users selesai")

	// ==================
	// CATEGORIES
	// ==================
	categories := []model.Category{
		{Name: "Makanan"},
		{Name: "Minuman"},
		{Name: "Snack"},
		{Name: "Dessert"},
	}
	for i := range categories {
		config.DB.FirstOrCreate(&categories[i], model.Category{Name: categories[i].Name})
	}
	fmt.Println("✅ Categories selesai")

	// ==================
	// BRANDS
	// ==================
	brands := []model.Brand{
		{Name: "Homemade"},
		{Name: "Indofood"},
		{Name: "Nestle"},
		{Name: "Local"},
	}
	for i := range brands {
		config.DB.FirstOrCreate(&brands[i], model.Brand{Name: brands[i].Name})
	}
	fmt.Println("✅ Brands selesai")

	// ==================
	// UNITS
	// ==================
	units := []model.Unit{
		{Name: "Pcs"},
		{Name: "Porsi"},
		{Name: "Gelas"},
		{Name: "Kg"},
		{Name: "Tray"},
	}
	for i := range units {
		config.DB.FirstOrCreate(&units[i], model.Unit{Name: units[i].Name})
	}
	fmt.Println("✅ Units selesai")

	// ==================
	// VARIANT CATEGORIES
	// ==================
	type OptionSeed struct {
		Name            string
		AdditionalPrice float64
	}
	type VariantSeed struct {
		Name      string
		MinSelect int
		MaxSelect int
		Options   []OptionSeed
	}

	variantSeeds := []VariantSeed{
		{
			Name: "Level Pedas", MinSelect: 1, MaxSelect: 1,
			Options: []OptionSeed{
				{"Level 1", 0}, {"Level 3", 0}, {"Level 5", 2000},
			},
		},
		{
			Name: "Topping", MinSelect: 0, MaxSelect: 3,
			Options: []OptionSeed{
				{"Keju", 3000}, {"Telur", 2000}, {"Sosis", 4000},
			},
		},
		{
			Name: "Suhu Minuman", MinSelect: 1, MaxSelect: 1,
			Options: []OptionSeed{
				{"Panas", 0}, {"Dingin", 1000},
			},
		},
		{
			Name: "Ukuran", MinSelect: 1, MaxSelect: 1,
			Options: []OptionSeed{
				{"Regular", 0}, {"Large", 5000},
			},
		},
	}

	variantMap := map[string]model.VariantCategory{}
	for _, vs := range variantSeeds {
		var existing model.VariantCategory
		err := config.DB.Where("name = ?", vs.Name).First(&existing).Error
		if err != nil {
			var options []model.VariantOption
			for _, o := range vs.Options {
				options = append(options, model.VariantOption{
					Name:            o.Name,
					AdditionalPrice: o.AdditionalPrice,
				})
			}
			newVariant := model.VariantCategory{
				Name:      vs.Name,
				MinSelect: vs.MinSelect,
				MaxSelect: vs.MaxSelect,
				Status:    "active",
				Options:   options,
			}
			config.DB.Create(&newVariant)
			variantMap[vs.Name] = newVariant
		} else {
			variantMap[vs.Name] = existing
		}
	}
	fmt.Println("✅ Variants selesai")

	// ==================
	// PRODUCTS
	// ==================
	type ProductSeed struct {
		Name         string
		CatIdx       int
		BrandIdx     int
		UnitIdx      int
		Price        float64
		VariantNames []string
	}

	productSeeds := []ProductSeed{
		{"Ayam Geprek", 0, 0, 1, 15000, []string{"Level Pedas", "Topping"}},
		{"Nasi Goreng", 0, 0, 1, 18000, []string{"Level Pedas"}},
		{"Mie Goreng", 0, 1, 1, 12000, []string{"Level Pedas"}},
		{"Soto Ayam", 0, 0, 1, 16000, []string{}},
		{"Es Teh Manis", 1, 0, 2, 5000, []string{"Ukuran"}},
		{"Kopi Susu", 1, 0, 2, 12000, []string{"Suhu Minuman", "Ukuran"}},
		{"Jus Alpukat", 1, 0, 2, 15000, []string{"Ukuran"}},
		{"Es Jeruk", 1, 0, 2, 8000, []string{"Ukuran"}},
		{"Keripik Singkong", 2, 3, 0, 8000, []string{}},
		{"Pisang Goreng", 2, 0, 0, 10000, []string{"Topping"}},
		{"Es Krim", 3, 2, 0, 13000, []string{"Topping"}},
		{"Pudding Coklat", 3, 0, 0, 11000, []string{}},
	}

	var createdProducts []model.Product
	for _, p := range productSeeds {
		var existing model.Product
		err := config.DB.Where("name = ?", p.Name).First(&existing).Error
		if err != nil {
			product := model.Product{
				Name:       p.Name,
				CategoryID: &categories[p.CatIdx].ID,
				BrandID:    &brands[p.BrandIdx].ID,
				UnitID:     &units[p.UnitIdx].ID,
				Price:      p.Price,
				Status:     "active",
			}
			config.DB.Create(&product)

			if len(p.VariantNames) > 0 {
				var variantObjs []model.VariantCategory
				for _, vname := range p.VariantNames {
					if v, ok := variantMap[vname]; ok {
						variantObjs = append(variantObjs, v)
					}
				}
				if len(variantObjs) > 0 {
					config.DB.Model(&product).Association("Variants").Replace(variantObjs)
				}
			}
			createdProducts = append(createdProducts, product)
		} else {
			createdProducts = append(createdProducts, existing)
		}
	}
	fmt.Println("✅ Products selesai")

	// ==================
	// PAYMENT METHODS
	// ==================
	payments := []model.PaymentMethod{
		{Name: "Cash"},
		{Name: "QRIS"},
		{Name: "Transfer Bank"},
	}
	for i := range payments {
		config.DB.FirstOrCreate(&payments[i], model.PaymentMethod{Name: payments[i].Name})
	}
	fmt.Println("✅ Payment Methods selesai")

	// ==================
// PROMOS
// ==================
type PromoSeed struct {
	Name        string
	PromoType   string
	MinTotal    float64
	DiscountPct float64
	CutPrice    float64
	VoucherCode string
	MaxUsage    int
	UsedCount   int
	Status      string
	StartDate   time.Time
	EndDate     time.Time
}

	promoSeeds := []PromoSeed{
		{
			Name:        "Diskon 10%",
			PromoType:   "discount",
			MinTotal:    50000,
			DiscountPct: 10,
			VoucherCode: "DISC10",
			MaxUsage:    100,
			UsedCount:   0,
			Status:      "active",
			StartDate:   time.Date(2026, 3, 1, 0, 0, 0, 0, time.Local),
			EndDate:     time.Date(2026, 3, 31, 23, 59, 59, 0, time.Local),
		},
		{
			Name:        "Hemat 5 Ribu",
			PromoType:   "cut_price",
			MinTotal:    30000,
			CutPrice:    5000,
			VoucherCode: "HEMAT5K",
			MaxUsage:    50,
			UsedCount:   0,
			Status:      "active",
			StartDate:   time.Date(2026, 3, 1, 0, 0, 0, 0, time.Local),
			EndDate:     time.Date(2026, 3, 31, 23, 59, 59, 0, time.Local),
		},
		{
			Name:        "Promo Lebaran",
			PromoType:   "discount",
			MinTotal:    100000,
			DiscountPct: 20,
			VoucherCode: "LEBARAN20",
			MaxUsage:    30,
			UsedCount:   0,
			Status:      "inactive",
			StartDate:   time.Date(2026, 4, 1, 0, 0, 0, 0, time.Local),
			EndDate:     time.Date(2026, 4, 10, 23, 59, 59, 0, time.Local),
		},
	}

	for _, p := range promoSeeds {
		var existing model.Promo
		if config.DB.Where("voucher_code = ?", p.VoucherCode).First(&existing).Error != nil {

			promo := model.Promo{
				Name:        p.Name,
				PromoType:   p.PromoType,
				AppliesTo:   "all",
				Condition:   "total",
				MinTotal:    p.MinTotal,
				DiscountPct: p.DiscountPct,
				CutPrice:    p.CutPrice,
				VoucherType: "custom",
				VoucherCode: p.VoucherCode,
				MaxUsage:    p.MaxUsage,
				UsedCount:   p.UsedCount,
				Status:      p.Status,
				StartDate:   p.StartDate,
				EndDate:     p.EndDate,
				ActiveDays:  "1,2,3,4,5,6,7",
				StartTime:   "00:00",
				EndTime:     "23:59",
			}

			config.DB.Create(&promo)
		}
	}

	fmt.Println("✅ Promos selesai")

	// ==================
	// SALES (Sample Transaksi)
	// ==================
	var kasir model.User
	config.DB.Where("username = ?", "budi").First(&kasir)

	type ItemSeed struct {
		ProductIdx int
		Qty        int
	}

	sampleSales := []struct {
		PaymentIdx int
		Items      []ItemSeed
	}{
		{0, []ItemSeed{{0, 2}, {4, 1}}},
		{1, []ItemSeed{{1, 1}, {5, 2}}},
		{0, []ItemSeed{{2, 3}, {6, 1}}},
		{2, []ItemSeed{{3, 1}, {7, 2}, {8, 1}}},
		{1, []ItemSeed{{4, 4}, {9, 2}}},
		{0, []ItemSeed{{0, 1}, {5, 1}, {10, 1}}},
		{1, []ItemSeed{{1, 2}, {11, 3}}},
		{0, []ItemSeed{{6, 2}, {7, 1}}},
	}

	for i, s := range sampleSales {
		invoiceNumber := fmt.Sprintf("INV-20260302-%04d", i+1)

		var existing model.Sale
		if config.DB.Where("invoice_number = ?", invoiceNumber).First(&existing).Error != nil {
			var saleItems []model.SaleItem
			var subtotal float64

			for _, item := range s.Items {
				if item.ProductIdx >= len(createdProducts) {
					continue
				}
				p := createdProducts[item.ProductIdx]
				itemSubtotal := p.Price * float64(item.Qty)
				subtotal += itemSubtotal
				saleItems = append(saleItems, model.SaleItem{
					ProductID: p.ID,
					Quantity:  item.Qty,
					BasePrice: p.Price,
					Subtotal:  itemSubtotal,
				})
			}

			sale := model.Sale{
				InvoiceNumber:   invoiceNumber,
				UserID:          kasir.ID,
				PaymentMethodID: payments[s.PaymentIdx].ID,
				Subtotal:        subtotal,
				DiscountTotal:   0,
				GrandTotal:      subtotal,
				Items:           saleItems,
			}
			config.DB.Create(&sale)
		}
	}
	fmt.Println("✅ Sales selesai")

	fmt.Println("")
	fmt.Println("🎉 Seed data berhasil!")
	fmt.Println("================================")
	fmt.Println("📌 Akun login yang tersedia:")
	fmt.Println("   Username: admin    | Password: password123 | Role: Owner")
	fmt.Println("   Username: budi     | Password: password123 | Role: Kasir")
	fmt.Println("   Username: siti     | Password: password123 | Role: Admin")
	fmt.Println("   Username: andi     | Password: password123 | Role: Supervisor")
	fmt.Println("================================")
}