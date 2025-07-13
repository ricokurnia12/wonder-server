package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

	// Hitung total hasil setelah filter
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count events"})
		return
	}

	sortBy := strings.ToLower(strings.TrimSpace(c.DefaultQuery("sortBy", "date")))
	sortOrder := strings.ToLower(strings.TrimSpace(c.DefaultQuery("sortOrder", "asc")))

	// --- DEBUGGING START ---
	fmt.Printf("DEBUG - Raw sortBy: '%s', Processed sortBy: '%s'\n", c.Query("sortBy"), sortBy)
	fmt.Printf("DEBUG - Raw sortOrder: '%s', Processed sortOrder: '%s'\n", c.Query("sortOrder"), sortOrder)
	// --- DEBUGGING END ---

	// Default sorting
	orderClause := "date ASC"

	switch sortBy {
	case "date":
		if sortOrder == "desc" {
			orderClause = "date DESC"
		} else {
			orderClause = "date ASC"
		}
	case "title":
		if sortOrder == "desc" {
			orderClause = "title DESC"
		} else {
			orderClause = "title ASC"
		}
	}

	// --- CRITICAL DEBUGGING POINT ---
	// This will show what value `orderClause` holds just before it's used by GORM.
	fmt.Printf("DEBUG - ***Final orderClause passed to GORM: '%s'***\n", orderClause)
	// --- END CRITICAL DEBUGGING POINT ---

	if err := query.
		Limit(limit).
		Offset(offset).
		Order(orderClause). // This is where the orderClause is applied
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
		"sortBy":     sortBy,
		"sortOrder":  sortOrder,
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
	var total int64

	// Query params
	search := c.Query("search")
	category := c.Query("category")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	offset := (page - 1) * limit

	// Query building
	query := database.DB.Model(&models.Event{})

	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	query.Count(&total)
	err := query.Limit(limit).Offset(offset).Order("date asc").Find(&events).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  events,
		"total": total,
		"page":  page,
		"limit": limit,
	})
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

// GetEventByID handles GET /events/:id
func GetEventByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var event models.Event
	if err := database.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// UpdateEvent handles PUT /events/:id
func UpdateEvent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var event models.Event
	if err := database.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	var input models.Event
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update hanya field yang diperlukan
	event.Title = input.Title
	event.Slug = input.Slug
	event.Description = input.Description
	event.Content = input.Content
	event.EnglishContent = input.EnglishContent
	event.Date = input.Date
	event.Location = input.Location
	event.Province = input.Province
	event.Category = input.Category
	event.Image = input.Image
	event.Featured = input.Featured

	if err := database.DB.Save(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, event)
}
