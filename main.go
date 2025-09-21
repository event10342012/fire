package main

import (
	"fire/internal/web/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	server := InitWebserver()
	server.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})
	err := server.Run(":8080")
	if err != nil {
		return
	}
}

func useJwt(server *gin.Engine) {
	login := &middleware.LoginJwtMiddlewareBuilder{}
	server.Use(login.CheckLogin())
}

//func useSession(server *gin.Engine) {
//	login := &middleware.LoginMiddlewareBuilder{}
//	// 存储数据的，也就是你 userId 存哪里
//	// 直接存 cookie
//	store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
//		[]byte("6jfbF1G0D2WcRjAZRq3Y2K47AGdL9nWT"),
//		[]byte("6jfbF1G0D2WcRjAZRq3Y2K47AGdL9nWS"))
//
//	if err != nil {
//		panic(err)
//	}
//	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
//}
