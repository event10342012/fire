package middleware

import (
	"encoding/gob"
	"fire/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginJwtMiddlewareBuilder struct {
}

func (m *LoginJwtMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	gob.Register(time.Time{})
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			// 不需要登录校验
			return
		}

		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		segs := strings.Split(authCode, " ")
		if len(segs) != 2 || segs[0] != "Bearer" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := segs[1]
		var uc web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return []byte(web.JwtKey), nil
		})

		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if uc.UserAgent != ctx.GetHeader("User-Agent") {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// refresh expiration time
		expireTime := uc.ExpiresAt
		if expireTime.Sub(time.Now()) < time.Hour*24 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7))
			tokenStr, err = token.SignedString([]byte(web.JwtKey))
			ctx.Header("x-jwt-token", tokenStr)
			if err != nil {
				log.Println(err)
			}
		}
		ctx.Set("user", uc)
	}
}
