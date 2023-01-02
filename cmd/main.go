package main

import (
	"github.com/gin-gonic/gin"

	"fire/src"
)

var router = gin.Default()

func main() {
	router.LoadHTMLGlob("templates/*")

	auth.AddAuthRoutes(router)

	router.Run(":8080")
}
