package main

import (
	"fire/auth"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/login", auth.LoginHandler)
	r.GET("/callback", auth.CallbackHandler)

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
