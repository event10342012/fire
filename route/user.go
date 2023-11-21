package route

import (
	"fire/controller"
	"github.com/gin-gonic/gin"
)

func SetupUserRoute(r *gin.RouterGroup) {
	userRoute := r.Group("user")

	userRoute.GET("/", controller.GetUsers)
	userRoute.GET("/:id", controller.GetUser)
	userRoute.POST("/", controller.CreateUser)
	userRoute.PUT("/:id", controller.UpdateUser)
	userRoute.DELETE("/:id", controller.DeleteUser)
}
