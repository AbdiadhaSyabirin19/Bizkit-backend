package handler

import (
	"bizkit-backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetAllCustomers(c *gin.Context) {
	search := c.Query("search")
	customers, err := service.GetAllCustomers(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": customers})
}

func GetCustomerByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}
	customer, err := service.GetCustomerByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": customer})
}

func CreateCustomer(c *gin.Context) {
	var req service.CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Input tidak valid", "error": err.Error()})
		return
	}
	customer, err := service.CreateCustomer(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Pelanggan berhasil ditambahkan", "data": customer})
}

func UpdateCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}
	var req service.CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Input tidak valid", "error": err.Error()})
		return
	}
	customer, err := service.UpdateCustomer(uint(id), req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data pelanggan diperbarui", "data": customer})
}

func DeleteCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}
	if err := service.DeleteCustomer(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghapus data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pelanggan berhasil dihapus"})
}
