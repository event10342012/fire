package web

import (
	"errors"
	"fire/internal/domain"
	"fire/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`

	JwtKey = "6jfbF1G0D2WcRjAZRq3Y2K47AGdL9nWT"
)

type UserHandler struct {
	emailRegexPattern    *regexp.Regexp
	passwordRegexPattern *regexp.Regexp
	svc                  *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRegexPattern:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexPattern: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:                  svc,
	}
}

func (handler *UserHandler) RegisterRoutes(server *gin.Engine) {
	userGroup := server.Group("/users")
	userGroup.POST("/login", handler.LoginJwt)
	userGroup.POST("/signup", handler.Signup)
	userGroup.GET("/profile", handler.Profile)
	userGroup.POST("/edit", handler.Edit)
}

func (handler *UserHandler) Login(ctx *gin.Context) {
	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req loginReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}

	user, err := handler.svc.Login(ctx, req.Email, req.Password)
	switch {
	case err == nil:
		sess := sessions.Default(ctx)
		sess.Set("userId", user.ID)
		sess.Options(sessions.Options{
			// 十五分钟
			MaxAge: 900,
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "System error")
			return
		}
		ctx.String(http.StatusOK, "Login success")
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		ctx.String(http.StatusOK, "email or password is invalid")
	default:
		ctx.String(http.StatusOK, "System error")
	}
}

func (handler *UserHandler) LoginJwt(ctx *gin.Context) {
	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req loginReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}

	user, err := handler.svc.Login(ctx, req.Email, req.Password)
	switch {
	case err == nil:
		uc := UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			},
			UserID: user.ID,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
		tokenString, err := token.SignedString([]byte(JwtKey))
		if err != nil {
			ctx.String(http.StatusOK, "System error")
			return
		}
		ctx.Header("x-jwt-token", tokenString)
		ctx.String(http.StatusOK, "Login success")
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		ctx.String(http.StatusOK, "email or password is invalid")
	default:
		ctx.String(http.StatusOK, "System error")
	}
}

func (handler *UserHandler) Signup(ctx *gin.Context) {
	type signupReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req signupReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}

	isEmail, err := handler.emailRegexPattern.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "email is invalid")
		return
	}

	isPassword, err := handler.passwordRegexPattern.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "system error")
	}
	if !isPassword {
		ctx.String(http.StatusOK, "password is invalid")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "password is not match")
		return
	}

	err = handler.svc.Signup(ctx, domain.User{
		Email:       req.Email,
		Password:    req.Password,
		IsSuperUser: false,
		IsActive:    true,
	})
	switch {
	case err == nil:
		ctx.String(http.StatusOK, "Signup success")
	case errors.Is(err, service.ErrDuplicateEmail):
		ctx.String(http.StatusOK, "Email is already exist")
	default:
		ctx.String(http.StatusOK, "Signup failed")
	}
}

func (handler *UserHandler) Profile(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	uid, ok := sess.Get("userId").(int64)
	if !ok {
		ctx.String(http.StatusOK, "System error")
		return
	}

	u, err := handler.svc.FindById(ctx, int(uid))
	if err != nil {
		ctx.String(http.StatusOK, "System error")
		return
	}

	type User struct {
		Nickname string    `json:"nickname"`
		Email    string    `json:"email"`
		AboutMe  string    `json:"aboutMe"`
		Birthday time.Time `json:"birthday"`
	}

	ctx.JSON(http.StatusOK, User{
		Nickname: u.Nickname,
		Email:    u.Email,
		AboutMe:  u.AboutMe,
		Birthday: u.Birthday,
	})
}

func (handler *UserHandler) Edit(ctx *gin.Context) {
	type editReq struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req editReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}

	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "birthday is invalid")
		return
	}

	sess := sessions.Default(ctx)
	uid, ok := sess.Get("userId").(int64)
	if !ok {
		ctx.String(http.StatusOK, "System error")
		return
	}

	err = handler.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		ID:       uid,
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})

	if err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}

	ctx.String(http.StatusOK, "Edit success")
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"userId"`
}
