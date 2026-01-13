//go:build wireinject

package startup

import (
	"fire/internal/repository"
	"fire/internal/repository/cache"
	"fire/internal/repository/dao"
	"fire/internal/service"
	"fire/internal/web"
	ijwt "fire/internal/web/jwt"
	"fire/ioc"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebserver() *gin.Engine {
	wire.Build(InitDB, ioc.InitRedis,
		dao.NewUserDAO,
		cache.NewCodeCache, cache.NewUserCache,
		repository.NewUserRepository, repository.NewCodeRepository,
		ioc.InitSMS, ioc.InitGoogleService, service.NewCodeService, service.NewUserService,
		web.NewUserHandler, ijwt.NewRedisJWTHandler, web.NewOAuth2GoogleHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer)
	return gin.Default()
}
