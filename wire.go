//go:build wireinject

package main

import (
	"fire/internal/repository"
	"fire/internal/repository/cache"
	"fire/internal/repository/dao"
	"fire/internal/service"
	"fire/internal/web"
	"fire/ioc"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebserver() *gin.Engine {
	wire.Build(ioc.InitDB, ioc.InitRedis,
		dao.NewUserDAO,
		cache.NewCodeCache, cache.NewUserCache,
		repository.NewUserRepository, repository.NewCodeRepository,
		ioc.InitSMS, service.NewCodeService, service.NewUserService,
		web.NewUserHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer)
	return gin.Default()
}
