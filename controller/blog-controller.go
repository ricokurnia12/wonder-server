package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ricokurnia12/wonder-server/database"
	"github.com/ricokurnia12/wonder-server/models"
)

func GetBlogPosts(c *gin.Context) {
	var posts []models.BlogPost
	database.DB.Preload("Author").Find(&posts)
	c.JSON(http.StatusOK, posts)
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
