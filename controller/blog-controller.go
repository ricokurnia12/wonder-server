package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ricokurnia12/wonder-server/database"
	"github.com/ricokurnia12/wonder-server/models"
)

func GetBlogPosts(c *gin.Context) {
	var posts []models.BlogPost
	var total int64

	// Query params
	search := c.Query("search")
	category := c.Query("category")
	sort := c.DefaultQuery("sort", "date_desc")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	featured := c.DefaultQuery("featured", "") // <-- fix

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	offset := (page - 1) * limit

	// Base query
	query := database.DB.Model(&models.BlogPost{})

	// Filter: search, category, featured
	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}
	if category != "" {
		query = query.Where("category ILIKE ?", "%"+category+"%")
	}
	if featured == "true" {
		query = query.Where("featured = ?", true)
	}

	// Hitung total
	query.Count(&total)

	// Sorting
	switch sort {
	case "featured_desc":
		query = query.Order("featured DESC").Order("date DESC")
	case "featured_asc":
		query = query.Order("featured ASC").Order("date DESC")
	case "date_asc":
		query = query.Order("date ASC")
	case "date_desc":
		fallthrough
	default:
		query = query.Order("date DESC")
	}

	// Fetch
	err := query.Limit(limit).Offset(offset).Find(&posts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  posts,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func GetBlogPostById(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var post models.BlogPost
	if err := database.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

func GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")

	var post models.BlogPost
	if err := database.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func UpdateBlogPost(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Struct khusus untuk update parsial
	var body struct {
		Title          *string `json:"title"`
		EnglishTitle   *string `json:"english_title"`
		Slug           *string `json:"slug"`
		Excerpt        *string `json:"excerpt"`
		EnglishExcerpt *string `json:"english_excerpt"`
		Content        *string `json:"content"`
		EnglishContent *string `json:"english_content"`
		Date           *string `json:"date"`
		ReadTime       *int    `json:"readTime"`
		Category       *string `json:"category"`
		CoverImage     *string `json:"coverImage"`
		Featured       *bool   `json:"featured"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cari post yang akan diupdate
	var existingPost models.BlogPost
	if err := database.DB.First(&existingPost, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Validasi jika ingin set Featured jadi true
	if body.Featured != nil && *body.Featured && !existingPost.Featured {
		var featuredCount int64
		if err := database.DB.Model(&models.BlogPost{}).Where("featured = ?", true).Count(&featuredCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count featured posts"})
			return
		}
		if featuredCount >= 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum of 3 featured posts allowed"})
			return
		}
	}

	// Bangun map untuk update data yang hanya dikirim
	updateData := make(map[string]interface{})
	if body.Title != nil {
		updateData["title"] = *body.Title
	}
	if body.EnglishTitle != nil {
		updateData["english_title"] = *body.EnglishTitle
	}
	if body.Slug != nil {
		updateData["slug"] = *body.Slug
	}
	if body.Excerpt != nil {
		updateData["excerpt"] = *body.Excerpt
	}
	if body.EnglishExcerpt != nil {
		updateData["english_excerpt"] = *body.EnglishExcerpt
	}
	if body.Content != nil {
		updateData["content"] = *body.Content
	}
	if body.EnglishContent != nil {
		updateData["english_content"] = *body.EnglishContent
	}
	if body.Date != nil {
		updateData["date"] = *body.Date
	}
	if body.ReadTime != nil {
		updateData["read_time"] = *body.ReadTime
	}
	if body.Category != nil {
		updateData["category"] = *body.Category
	}
	if body.CoverImage != nil {
		updateData["cover_image"] = *body.CoverImage
	}
	if body.Featured != nil {
		updateData["featured"] = *body.Featured
	}

	// Lakukan update jika ada data yang dikirim
	if len(updateData) > 0 {
		if err := database.DB.Model(&existingPost).Updates(updateData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

func CreateBlogPost(c *gin.Context) {
	var post models.BlogPost
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&post)
	c.JSON(http.StatusOK, post)
}

func DeleteBlogPost(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := database.DB.Delete(&models.BlogPost{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
