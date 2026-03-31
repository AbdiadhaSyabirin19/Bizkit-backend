package handler

import (
	"net/http"

	"bizkit-backend/internal/service"

	"github.com/gin-gonic/gin"
)

func GetSalesReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	onlyDiscounted := c.Query("only_discounted") == "true"

	result, err := service.GetSalesReport(startDate, endDate, onlyDiscounted)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": result})
}

func GetTrendReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	result, err := service.GetTrendReport(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": result})
}

func GetAttendanceReport(c *gin.Context) {
	date := c.Query("date")

	result, err := service.GetAttendanceReport(date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": result})
}

func GetShiftReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	result, err := service.GetShiftReport(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": result})
}

func GetDashboardSummary(c *gin.Context) {
	data, err := service.GetDashboardSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": data})
}
