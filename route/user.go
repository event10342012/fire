package route

import (
	"fire/controller"
	"github.com/gin-gonic/gin"
)

func SetupUserRoute(r *gin.RouterGroup) {
	user := r.Group("user")

	user.GET("/:id", controller.GetUserByID)
	user.POST("/", controller.PostUser)
	user.PUT("/:id", controller.PutUser)
}
