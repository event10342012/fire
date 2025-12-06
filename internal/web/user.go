package web

import (
	"errors"
	"fire/internal/domain"
	"fire/internal/repository"
	"fire/internal/service"
	"net/http"
	"time"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`

	bizLogin = "login"
)

type UserHandler struct {
	emailRegexPattern    *regexp.Regexp
	passwordRegexPattern *regexp.Regexp
	userSvc              service.UserService
	codeSvc              service.CodeService
	jwtHandler
}

func NewUserHandler(userSvc service.UserService, codeSvc service.CodeService) *UserHandler {
	return &UserHandler{
		emailRegexPattern:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexPattern: regexp.MustCompile(passwordRegexPattern, regexp.None),
		userSvc:              userSvc,
		codeSvc:              codeSvc,
	}
}

func (handler *UserHandler) RegisterRoutes(server *gin.Engine) {
	userGroup := server.Group("/users")
	userGroup.POST("/login", handler.LoginJwt)
	userGroup.POST("/signup", handler.Signup)
	userGroup.GET("/profile", handler.Profile)
	userGroup.POST("/edit", handler.Edit)
	userGroup.POST("/refresh_token", handler.RefreshToken)

	userGroup.POST("/login_sms/code/send", handler.SendSMSLoginCode)
	userGroup.POST("login_sms", handler.LoginSMS)
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

	user, err := handler.userSvc.Login(ctx, req.Email, req.Password)
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

	user, err := handler.userSvc.Login(ctx, req.Email, req.Password)
	switch {
	case err == nil:
		handler.setJWTToken(ctx, user.ID)
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

	err = handler.userSvc.Signup(ctx, domain.User{
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

	u, err := handler.userSvc.FindById(ctx, uid)
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

	err = handler.userSvc.UpdateNonSensitiveInfo(ctx, domain.User{
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

func (handler *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		return
	}

	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "phone is empty",
		})
		return
	}

	err = handler.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, Result{
			Code: 0,
			Msg:  "send success",
		})
	case errors.Is(err, service.ErrCodeSendTooMany):
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "send too many",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "System error",
		})
	}
}

func (handler *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}

	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		return
	}

	ok, err := handler.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "System error",
		})
		return
	}

	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "code is invalid",
		})
		return
	}

	u, err := handler.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "System error",
		})
	}
	handler.setJWTToken(ctx, u.ID)
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "Login success",
	})
}

func (handler *UserHandler) FindOrCreate(ctx *gin.Context, phone string) (domain.User, error) {
	u, err := handler.userSvc.FindByPhone(ctx, phone)
	if !errors.Is(err, repository.ErrUserNotFound) {
		return u, err
	}

	err = handler.userSvc.Create(ctx, domain.User{
		Phone: phone,
	})

	if err != nil && !errors.Is(err, repository.ErrDuplicateUser) {
		return domain.User{}, err
	}
	// may cause error due to replication delay between master and slave db
	return handler.userSvc.FindByPhone(ctx, phone)
}

func (handler *UserHandler) RefreshToken(ctx *gin.Context) {
	tokenStr := ExtractToken(ctx)
	var rc RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	handler.setJWTToken(ctx, rc.UserID)
	ctx.Status(http.StatusOK)
}
