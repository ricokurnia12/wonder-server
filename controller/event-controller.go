package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ricokurnia12/wonder-server/database"
	"github.com/ricokurnia12/wonder-server/models"
	"gorm.io/gorm"
)

func GetEventsClient(c *gin.Context) {
	var events []models.Event
	var total int64

	db := database.DB.Debug() // Debug untuk lihat query di console, hapus di production

	// Default pagination
	page := 1
	limit := 10

	// Ambil query param page
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
		if page < 1 {
			page = 1
		}
	}

	// Ambil query param limit
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
		if limit < 1 {
			limit = 10
		}
	}

	offset := (page - 1) * limit

	// Buat query awal
	query := db.Model(&models.Event{})

	// Filter berdasarkan title
	if title := c.Query("title"); title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	// Filter berdasarkan category
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	// Filter berdasarkan startDate
	if startDateStr := c.Query("startDate"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startDate format. Use YYYY-MM-DD."})
			return
		}
		query = query.Where("date >= ?", startDate)
	}

	// Filter berdasarkan endDate
	if endDateStr := c.Query("endDate"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endDate format. Use YYYY-MM-DD."})
			return
		}
		// Tambahkan 1 hari dikurangi 1 detik (jadi jam 23:59:59)
		endOfDay := endDate.Add(24*time.Hour - time.Second)
		query = query.Where("date <= ?", endOfDay)
	}

	// Hitung total hasil setelah filter
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count events"})
		return
	}

	// Ambil data dengan limit dan offset
	if err := query.
		Limit(limit).
		Offset(offset).
		Order("date ASC").
		Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	// Berikan response
	c.JSON(http.StatusOK, gin.H{
		"data":       events,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": int((total + int64(limit) - 1) / int64(limit)),
	})
}
func GetEventBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug is required"})
		return
	}

	var event models.Event
	if err := database.DB.
		Where("slug = ?", slug).
		First(&event).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch event"})
		}
		return
	}

	c.JSON(http.StatusOK, event)

}

func GetEvents(c *gin.Context) {
	var events []models.Event
	database.DB.Find(&events)
	c.JSON(http.StatusOK, events)
}

func CreateEvent(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&event)
	c.JSON(http.StatusOK, event)
}
