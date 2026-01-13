package ioc

import (
	"fire/internal/web"
	ijwt "fire/internal/web/jwt"
	"fire/internal/web/middleware"
	"fire/pkg/ginx/middleware/ratelimit"
	"fire/pkg/limiter"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func InitWebServer(
	mdls []gin.HandlerFunc,
	userHdl *web.UserHandler,
	googleHdl *web.OAuth2GoogleHandler,
	artHdl *web.ArticleHandler,
) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	googleHdl.RegisterRoutes(server)
	artHdl.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(redisClient redis.Cmdable, handler ijwt.Handler) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			AllowCredentials: true,
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			ExposeHeaders:    []string{"x-jwt-token"},
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					//if strings.Contains(origin, "localhost") {
					return true
				}
				return strings.Contains(origin, "your_company.com")
			},
			MaxAge: 12 * time.Hour,
		}),
		ratelimit.NewBuilder(limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 1000)).Build(),
		(&middleware.LoginJwtMiddlewareBuilder{Handler: handler}).CheckLogin(),
		middleware.NewLoginJwtMiddlewareBuilder(handler).CheckLogin(),
	}
}
