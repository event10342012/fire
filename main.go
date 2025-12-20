package main

import (
	"fire/internal/web/middleware"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	initViper()
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

func initViper() {
	cfile := pflag.StringP("config", "c", "config/config.yml", "config file path")
	pflag.Parse()

	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cfile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	log.Printf("Load config file %s", viper.ConfigFileUsed())
}
