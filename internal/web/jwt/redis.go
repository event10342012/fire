package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var JWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK")
var RCJWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgA")

type RedisJWTHandler struct {
	signedMethod jwt.SigningMethod
	refreshKey   []byte
	client       redis.Cmdable
	rcExpiration time.Duration
}

func NewRedisJWTHandler(client redis.Cmdable) Handler {
	return &RedisJWTHandler{
		signedMethod: jwt.SigningMethodHS256,
		refreshKey:   RCJWTKey,
		client:       client,
		rcExpiration: time.Hour * 24 * 7,
	}
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID    int64
	Ssid      string
	UserAgent string
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	UserID int64
	Ssid   string
}

func (handler *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
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

func (handler *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-refresh-token", "")
	ctx.Header("x-jwt-token", "")
	uc := ctx.MustGet("user").(UserClaims)
	return handler.client.Set(ctx, fmt.Sprintf("users:ssids:%s", uc.Ssid), "", handler.rcExpiration).Err()
}

func (handler *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := handler.SetRefreshToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	return handler.SetJWTToken(ctx, uid, ssid)
}

func (handler *RedisJWTHandler) SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	rc := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(handler.rcExpiration)),
		},
		UserID: uid,
		Ssid:   ssid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, rc)
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		ctx.String(http.StatusOK, "System error")
		return fmt.Errorf("refresh token error: %w", err)
	}
	ctx.Header("x-refresh-token", tokenString)
	return nil
}

func (handler *RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	uc := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserID:    uid,
		Ssid:      ssid,
		UserAgent: ctx.GetHeader("User-Agent"),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenString)
	return nil
}

func (handler *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	cnt, err := handler.client.Exists(ctx, fmt.Sprintf("users:ssids:%s", ssid)).Result()
	if err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("invalid token")
	}
	return nil
}
