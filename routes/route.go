package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ricokurnia12/wonder-server/controller"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/events", controllers.GetEvents)
		api.GET("/eventsclient", controllers.GetEventsClient)
		api.POST("/events", controllers.CreateEvent)

		api.GET("/blogposts", controllers.GetBlogPosts)
		api.POST("/blogposts", controllers.CreateBlogPost)

		api.POST("/photos", controllers.UploadPhoto)
		api.GET("/photos", controllers.GetPhotos)
		api.GET("/photos/:id", controllers.GetPhoto)
		api.DELETE("/photos/:id", controllers.DeletePhoto)
		api.GET("/photos/paginated", controllers.GetPhotosPaginated)

		// optionally: serve static files
		api.Static("/uploads", "./uploads")
	}
}
