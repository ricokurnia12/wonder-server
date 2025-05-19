package controllers

import (
	"net/http"
	"os"
	"path/filepath"

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

	c.JSON(http.StatusOK, photo)
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
