package web

import (
	"fmt"
	"net/http"
	"strings"
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

type RefreshClaims struct {
	jwt.RegisteredClaims
	UserID int64
}

func (handler *jwtHandler) setRefreshToken(ctx *gin.Context, uid int64) error {
	rc := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		},
		UserID: uid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, rc)
	tokenString, err := token.SignedString([]byte(JwtKey))
	if err != nil {
		ctx.String(http.StatusOK, "System error")
		return fmt.Errorf("refresh token error: %w", err)
	}
	ctx.Header("x-refresh-token", tokenString)
	return nil
}

func (handler *jwtHandler) setJWTToken(ctx *gin.Context, uid int64) {
	err := handler.setRefreshToken(ctx, uid)
	if err != nil {
		ctx.String(http.StatusOK, "System error")
	}
	uc := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserID:    uid,
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

func ExtractToken(ctx *gin.Context) string {
	authCode := ctx.GetHeader("Authorization")
	if authCode == "" {
		return authCode
	}

	segs := strings.Split(authCode, " ")
	if len(segs) != 2 {
		return ""
	}

	tokenStr := segs[1]
	return tokenStr
}
