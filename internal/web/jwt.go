package web

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const JwtKey = "6jfbF1G0D2WcRjAZRq3Y2K47AGdL9nWT"

type jwtHandler struct {
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID    int64
	UserAgent string
}

func (handler *jwtHandler) setJWTToken(ctx *gin.Context, id int64) {
	uc := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserID:    id,
		UserAgent: ctx.GetHeader("User-Agent"),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenString, err := token.SignedString([]byte(JwtKey))
	if err != nil {
		ctx.String(http.StatusOK, "System error")
		return
	}
	ctx.Header("x-jwt-token", tokenString)
}
