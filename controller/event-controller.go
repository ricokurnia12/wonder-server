package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/ricokurnia12/wonder-server/database"
	"github.com/ricokurnia12/wonder-server/models"
)

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
