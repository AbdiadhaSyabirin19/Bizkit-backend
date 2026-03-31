package handler

import (
	"net/http"
	"strconv"

	"bizkit-backend/internal/service"

	"github.com/gin-gonic/gin"
)

func CreateReceivablePayment(c *gin.Context) {
	var req service.ReceivablePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Input tidak valid", "error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)

	payment, err := service.CreateReceivablePayment(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pembayaran berhasil disimpan", "data": payment})
}

func UpdateReceivablePayment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}

	var req service.ReceivablePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Input tidak valid", "error": err.Error()})
		return
	}

	payment, err := service.UpdateReceivablePayment(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran berhasil diperbarui", "data": payment})
}

func DeleteReceivablePayment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}

	if err := service.DeleteReceivablePayment(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghapus pembayaran", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran berhasil dihapus"})
}

func GetPaymentsBySaleID(c *gin.Context) {
	saleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Sale ID tidak valid"})
		return
	}

	items, err := service.GetPaymentsBySaleID(uint(saleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil data history pembayaran: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

func GetUnpaidSales(c *gin.Context) {
	customerName := c.Query("customer_name")

	sales, err := service.GetUnpaidSales(customerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil sales yang belum lunas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sales})
}

func GetReceivablePaymentByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}

	payment, err := service.GetReceivablePaymentByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Pembayaran tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payment})
}

func FixMigrate(c *gin.Context) {
	err := service.FixMigrate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Migrate Success"})
}

func SyncSales(c *gin.Context) {
	err := service.SyncSales()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Semua data penjualan lama berhasil disinkronkan dengan histori pembayaran otomatis!"})
}
