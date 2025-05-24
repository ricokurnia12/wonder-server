package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ricokurnia12/wonder-server/controller"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// events
		api.GET("/events", controllers.GetEvents)
		api.GET("/eventsclient", controllers.GetEventsClient)
		api.POST("/events", controllers.CreateEvent)
		api.GET("/events/:slug", controllers.GetEventBySlug)
		api.GET("/detail-to-edit/:id", controllers.GetEventByID)
		api.PUT("/events/:id", controllers.UpdateEvent)
		// blog
		api.GET("/blogposts", controllers.GetBlogPosts)
		api.POST("/blogposts", controllers.CreateBlogPost)
		api.GET("/blogposts/:id", controllers.GetBlogPostById)
		api.PUT("/blogposts/:id", controllers.UpdateBlogPost)
		api.DELETE("/blogposts/:id", controllers.DeleteBlogPost)
		// photos
		api.POST("/photos", controllers.UploadPhoto)
		api.GET("/photos", controllers.GetPhotos)
		api.GET("/photos/:id", controllers.GetPhoto)
		api.DELETE("/photos/:id", controllers.DeletePhoto)
		api.GET("/photos/paginated", controllers.GetPhotosPaginated)

		// optionally: serve static files
		api.Static("/uploads", "./uploads")
	}
}
