package handler

import (
	"net/http"

	"bizkit-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// SyncOfflineSales godoc
// @Summary      Sinkronisasi transaksi offline kasir
// @Description  Mengirim batch transaksi yang disimpan secara offline di device klien (Android/iOS).
//
//	Setiap transaksi harus memiliki `offline_id` berformat UUID v4 yang unik.
//	Endpoint ini bersifat IDEMPOTENT: memanggil dengan data yang sama berkali-kali aman —
//	transaksi yang sudah ada akan dilewati (skipped), bukan di-insert ulang.
//	Maksimal 100 transaksi per sekali request.
//
// @Tags         Sales
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      service.SyncRequest  true  "Batch transaksi offline"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Router       /sales/sync [post]
func SyncOfflineSales(c *gin.Context) {
	var req service.SyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Format request tidak valid",
			"detail":  err.Error(),
		})
		return
	}

	if len(req.Transactions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Tidak ada transaksi yang dikirim",
		})
		return
	}

	userID := c.MustGet("user_id").(uint)

	results, summary, err := service.SyncOfflineSales(req.Transactions, userID)
	if err != nil {
		// Error ini hanya terjadi jika batch melebihi batas maksimal
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// HTTP 200 selalu dikembalikan — hasil per-item ada di dalam results.
	// Klien harus memeriksa field "status" tiap item (created/skipped/failed).
	c.JSON(http.StatusOK, gin.H{
		"message": "Sinkronisasi selesai",
		"summary": summary,
		"results": results,
	})
}
