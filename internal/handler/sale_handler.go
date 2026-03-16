package handler

import (
	"net/http"
	"strconv"

	"bizkit-backend/internal/service"

	"github.com/gin-gonic/gin"
)

func CreateSale(c *gin.Context) {
	var req service.SaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Input tidak valid"})
		return
	}

	userID := c.MustGet("user_id").(uint)

	sale, err := service.CreateSale(req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Transaksi berhasil dibuat",
		"data":    sale,
	})
}

func GetAllSales(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	sales, err := service.GetAllSales(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": sales})
}

func GetSaleByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}

	sale, err := service.GetSaleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": sale})
}

func GetDailySales(c *gin.Context) {
	date := c.Query("date")
	source := c.Query("source")

	result, err := service.GetDailySales(date, source)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": result})
}

func UpdateSale(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}

	var req service.SaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Input tidak valid"})
		return
	}

	sale, err := service.UpdateSale(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaksi berhasil diperbarui",
		"data":    sale,
	})
}

func DeleteSale(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}

	if err := service.DeleteSale(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghapus data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaksi berhasil dihapus"})
}