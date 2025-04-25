package main

import (
	"fire/config"
	"fire/internal/repository"
	"fire/internal/repository/cache"
	"fire/internal/repository/dao"
	"fire/internal/service"
	localSms "fire/internal/service/sms/local"
	"fire/internal/web"
	"fire/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})

	server := initWebServer()
	codeSvc := initCodeSvc(redisClient)
	initUserHdl(db, redisClient, codeSvc, server)
	err := server.Run(":8080")
	if err != nil {
		return
	}
}

func initUserHdl(db *gorm.DB, redisClient redis.Cmdable, codeSvc service.CodeService, server *gin.Engine) {
	userSvc := initUserSvc(db, redisClient)
	userHdl := web.NewUserHandler(userSvc, codeSvc)
	userHdl.RegisterRoutes(server)
}

func initUserSvc(db *gorm.DB, redisClient redis.Cmdable) service.UserService {
	userDao := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(redisClient)
	userRepo := repository.NewUserRepository(userDao, userCache)
	return service.NewUserService(userRepo)
}

func initCodeSvc(redisClient redis.Cmdable) service.CodeService {
	codeCache := cache.NewCodeCache(redisClient)
	codeRepo := repository.NewCodeRepository(codeCache)
	smsSvc := localSms.NewService()
	return service.NewCodeService(codeRepo, smsSvc)
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
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
	}))

	//redisClient := redis.NewClient(&redis.Options{
	//	Addr:     config.Config.Redis.Addr,
	//	Password: "",
	//	DB:       0,
	//})
	//
	//server.Use(ratelimit.NewBuilder(redisClient, time.Minute, 100).Build())

	useJwt(server)
	return server
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
