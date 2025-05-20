package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ricokurnia12/wonder-server/database"
	"github.com/ricokurnia12/wonder-server/routes"
)

func main() {
	database.ConnectionDb()
	r := gin.Default()

	// Allow all origins (public API)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // true tidak diizinkan jika AllowOrigins = ["*"]
		MaxAge:           12 * time.Hour,
	}))

	os.MkdirAll("uploads", os.ModePerm)
	routes.SetupRoutes(r)
	r.Run(":8080")
}
