package main

import (
	"fire/route"
	"fire/server"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	v1 := r.Group("v1")

	r.GET("/login", server.LoginHandler)
	r.GET("/callback", server.CallbackHandler)

	route.SetupUserRoute(v1)

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
