package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ricokurnia12/wonder-server/database"
	"github.com/ricokurnia12/wonder-server/routes"
)

func main() {
	database.ConnectionDb()
	r := gin.Default()
	routes.SetupRoutes(r)
	r.Run(":8080")
}
