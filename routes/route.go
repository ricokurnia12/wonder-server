package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ricokurnia12/wonder-server/controller"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/events", controllers.GetEvents)
		api.POST("/events", controllers.CreateEvent)

		api.GET("/blogposts", controllers.GetBlogPosts)
		api.POST("/blogposts", controllers.CreateBlogPost)
	}
}
