package route

import (
	"fire/controller"
	"github.com/gin-gonic/gin"
)

func SetupAccountingRoute(r *gin.RouterGroup) {
	accountRoute := r.Group("accounting")

	accountRoute.GET("/", controller.GetAccountingTransactions)
	accountRoute.GET("/:id", controller.GetAccountingTransaction)
	accountRoute.POST("/", controller.CreateAccountingTransaction)
	accountRoute.PUT("/:id", controller.UpdateAccountingTransaction)
	accountRoute.DELETE("/:id", controller.DeleteAccountingTransaction)
}
