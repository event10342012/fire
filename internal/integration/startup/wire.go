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

var thirdPartySet = wire.NewSet(
	ioc.InitDB, ioc.InitRedis, ioc.InitLogger,
)

func InitWebserver() *gin.Engine {
	wire.Build(
		thirdPartySet,
		// dao
		dao.NewUserDAO, dao.NewArticleGormDAO,
		// cache
		cache.NewCodeCache, cache.NewUserCache,
		// repository
		repository.NewUserRepository, repository.NewCodeRepository, repository.NewArticleRepository,
		// service
		ioc.InitSMS, ioc.InitGoogleService, service.NewCodeService, service.NewUserService, service.NewArticleService,
		// handler
		web.NewUserHandler, ijwt.NewRedisJWTHandler, web.NewOAuth2GoogleHandler, web.NewArticleHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(
		dao.NewArticleGormDAO,
		repository.NewArticleRepository,
		service.NewArticleService,
		thirdPartySet,
		web.NewArticleHandler,
	)
	return &web.ArticleHandler{}
}
