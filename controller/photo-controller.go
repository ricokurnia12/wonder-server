package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ricokurnia12/wonder-server/database"
	"github.com/ricokurnia12/wonder-server/models"
)

func UploadPhoto(c *gin.Context) {
	title := c.PostForm("title")
	file, err := c.FormFile("photo")
	print(title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Foto wajib diunggah", "detail": err.Error()})
		return
	}

	// Simpan file ke folder uploads/
	path := filepath.Join("uploads", file.Filename)
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan file"})
		return
	}

	photo := models.Photo{Title: title, FilePath: path}
	database.DB.Create(&photo)

	c.JSON(http.StatusOK, path)
}

func GetPhotos(c *gin.Context) {
	var photos []models.Photo
	database.DB.Find(&photos)
	c.JSON(http.StatusOK, photos)
}

func GetPhoto(c *gin.Context) {
	var photo models.Photo
	if err := database.DB.First(&photo, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Foto tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, photo)
}

func DeletePhoto(c *gin.Context) {
	var photo models.Photo
	if err := database.DB.First(&photo, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Foto tidak ditemukan"})
		return
	}

	// Hapus file dari sistem
	os.Remove(photo.FilePath)

	database.DB.Delete(&photo)
	c.JSON(http.StatusOK, gin.H{"message": "Foto dihapus"})
}

func GetPhotosPaginated(c *gin.Context) {
	// Ambil query param page dan limit
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	// Konversi string ke int
	pageNum, err1 := strconv.Atoi(page)
	limitNum, err2 := strconv.Atoi(limit)
	if err1 != nil || err2 != nil || pageNum <= 0 || limitNum <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter page dan limit harus berupa angka positif"})
		return
	}

	// Hitung offset
	offset := (pageNum - 1) * limitNum

	var photos []models.Photo
	var total int64

	// Hitung total data
	database.DB.Model(&models.Photo{}).Count(&total)

	// Ambil data dengan limit dan offset
	database.DB.Limit(limitNum).Offset(offset).Order("id desc").Find(&photos)

	c.JSON(http.StatusOK, gin.H{
		"data":       photos,
		"page":       pageNum,
		"limit":      limitNum,
		"total":      total,
		"totalPages": int((total + int64(limitNum) - 1) / int64(limitNum)),
	})
}
