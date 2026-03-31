package handler

import (
	"net/http"

	"bizkit-backend/internal/service"

	"github.com/gin-gonic/gin"
)

func GetDailyAdvancedReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	mode := c.Query("mode")
	yearlyBreakdown := c.Query("yearly_breakdown") // "monthly" or "weekly"

	if mode == "" {
		mode = "daily" // default
	}

	result, err := service.GetDailyAdvancedReport(startDate, endDate, mode, yearlyBreakdown)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": result})
}
