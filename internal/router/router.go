package router

import (
	"bizkit-backend/internal/handler"
	"bizkit-backend/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Static("/uploads", "./uploads")

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://bizkit-frontend.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
		api.HEAD("/ping", func(c *gin.Context) {
			c.Status(200)
		})

		// Auth (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", handler.Login)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/me", handler.GetMe)
			protected.PUT("/auth/change-password", handler.ChangePassword)

			// Category
			protected.GET("/categories", handler.GetAllCategories)
			protected.GET("/categories/:id", handler.GetCategoryByID)
			protected.POST("/categories", handler.CreateCategory)
			protected.PUT("/categories/:id", handler.UpdateCategory)
			protected.DELETE("/categories/:id", handler.DeleteCategory)

			// Brand
			protected.GET("/brands", handler.GetAllBrands)
			protected.GET("/brands/:id", handler.GetBrandByID)
			protected.POST("/brands", handler.CreateBrand)
			protected.PUT("/brands/:id", handler.UpdateBrand)
			protected.DELETE("/brands/:id", handler.DeleteBrand)

			// Unit
			protected.GET("/units", handler.GetAllUnits)
			protected.GET("/units/:id", handler.GetUnitByID)
			protected.POST("/units", handler.CreateUnit)
			protected.PUT("/units/:id", handler.UpdateUnit)
			protected.DELETE("/units/:id", handler.DeleteUnit)

			// Variant
			protected.GET("/variants", handler.GetAllVariantCategories)
			protected.GET("/variants/:id", handler.GetVariantCategoryByID)
			protected.POST("/variants", handler.CreateVariantCategory)
			protected.PUT("/variants/:id", handler.UpdateVariantCategory)
			protected.DELETE("/variants/:id", handler.DeleteVariantCategory)

			// Product
			protected.GET("/products", handler.GetAllProducts)
			protected.GET("/products/:id", handler.GetProductByID)
			protected.POST("/products", handler.CreateProduct)
			protected.PUT("/products/:id", handler.UpdateProduct)
			protected.DELETE("/products/:id", handler.DeleteProduct)
			protected.GET("/products/:id/prices", handler.GetProductPrices)

			// Role
			protected.GET("/roles", handler.GetAllRoles)
			protected.GET("/roles/:id", handler.GetRoleByID)
			protected.POST("/roles", handler.CreateRole)
			protected.PUT("/roles/:id", handler.UpdateRole)
			protected.DELETE("/roles/:id", handler.DeleteRole)

			// User
			protected.GET("/users", handler.GetAllUsers)
			protected.GET("/users/:id", handler.GetUserByID)
			protected.POST("/users", handler.CreateUser)
			protected.PUT("/users/:id", handler.UpdateUser)
			protected.DELETE("/users/:id", handler.DeleteUser)

			// Payment Method
			protected.GET("/payment-methods", handler.GetAllPaymentMethods)
			protected.GET("/payment-methods/:id", handler.GetPaymentMethodByID)
			protected.POST("/payment-methods", handler.CreatePaymentMethod)
			protected.PUT("/payment-methods/:id", handler.UpdatePaymentMethod)
			protected.DELETE("/payment-methods/:id", handler.DeletePaymentMethod)

			// Promo
			protected.GET("/promos", handler.GetAllPromos)
			protected.GET("/promos/:id", handler.GetPromoByID)
			protected.POST("/promos", handler.CreatePromo)
			protected.PUT("/promos/:id", handler.UpdatePromo)
			protected.DELETE("/promos/:id", handler.DeletePromo)
			protected.GET("/products/:id/promos", handler.GetPromosByProduct)
			protected.POST("/promos/check", handler.CheckAutoPromos)
			protected.POST("/promos/check-voucher", handler.CheckVoucher)

			// Sales — /daily harus SEBELUM /:id
			protected.POST("/sales", handler.CreateSale)
			protected.GET("/sales", handler.GetAllSales)
			protected.GET("/sales/daily", handler.GetDailySales)
			protected.GET("/sales/:id", handler.GetSaleByID)
			protected.PUT("/sales/:id", handler.UpdateSale)
			protected.DELETE("/sales/:id", handler.DeleteSale)

			// Shifts — di dalam protected (butuh auth)
			protected.POST("/shifts/open", handler.OpenShift)
			protected.POST("/shifts/close", handler.CloseShift)
			protected.GET("/shifts/active", handler.GetActiveShift)
			protected.GET("/shifts/history", handler.GetShiftHistory)
			protected.GET("/shifts/:id/summary", handler.GetShiftSummary)
			protected.PATCH("/shifts/:id/notes", handler.UpdateShiftNotes)

			// Laporan
			protected.GET("/reports/sales", handler.GetSalesReport)
			protected.GET("/reports/trend", handler.GetTrendReport)
			protected.GET("/reports/attendance", handler.GetAttendanceReport)
			protected.GET("/reports/shift", handler.GetShiftReport)
			protected.GET("/reports/dashboard-summary", handler.GetDashboardSummary)

			// Generic Upload
			protected.POST("/upload", handler.UploadFile)

			// Customer
			protected.GET("/customers", handler.GetAllCustomers)
			protected.GET("/customers/:id", handler.GetCustomerByID)
			protected.POST("/customers", handler.CreateCustomer)
			protected.PUT("/customers/:id", handler.UpdateCustomer)
			protected.DELETE("/customers/:id", handler.DeleteCustomer)

			// Setting
			protected.GET("/settings", handler.GetSetting)
			protected.PUT("/settings", handler.UpdateSetting)

			// Outlet
			protected.GET("/outlets", handler.GetAllOutlets)
			protected.GET("/outlets/:id", handler.GetOutletByID)
			protected.POST("/outlets", handler.CreateOutlet)
			protected.PUT("/outlets/:id", handler.UpdateOutlet)
			protected.DELETE("/outlets/:id", handler.DeleteOutlet)

			// Multi Harga
			protected.GET("/price-categories", handler.GetAllPriceCategories)
			protected.POST("/price-categories", handler.CreatePriceCategory)
			protected.PUT("/price-categories/:id", handler.UpdatePriceCategory)
			protected.DELETE("/price-categories/:id", handler.DeletePriceCategory)
			protected.GET("/price-categories/:id/products", handler.GetProductPricesByCategory)
			protected.POST("/price-categories/:id/products", handler.UpsertProductPrice)
			protected.DELETE("/price-categories/:id/products/:product_id", handler.DeleteProductPrice)
			protected.GET("/price-categories/:id", handler.GetPriceCategoryByID)

			// Attendance
			protected.POST("/attendances/checkin", handler.CheckIn)
			protected.POST("/attendances/:id/checkout", handler.CheckOut)
			protected.GET("/attendances/today", handler.GetTodayAttendance)
			protected.GET("/attendances/history", handler.GetAttendanceHistory)
			
		}
	}
	return r
}