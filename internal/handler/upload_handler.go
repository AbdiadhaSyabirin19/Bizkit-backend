package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "File wajib diupload"})
		return
	}

	// Tentukan subfolder (misal: products, users, etc)
	folder := c.DefaultPostForm("folder", "general")
	uploadDir := "uploads/" + folder
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat direktori upload"})
		return
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := uploadDir + "/" + filename

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menyimpan file"})
		return
	}

	// Return URL path
	c.JSON(http.StatusOK, gin.H{
		"message": "Upload berhasil",
		"url":     "/" + savePath,
	})
}
