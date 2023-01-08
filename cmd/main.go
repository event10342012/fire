package main

import (
	"fire/src"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func main() {
	router.LoadHTMLGlob("templates/*")

	auth.AddAuthRoutes(router)

	router.Run(":8080")
}
