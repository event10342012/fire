package web

import (
	"fire/internal/service"
	"fire/internal/service/oauth2/google"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type OAuth2GoogleHandler struct {
	googleAuthSvc   google.AuthService
	userSvc         service.UserService
	key             []byte
	stateCookieName string
	jwtHandler
}

func NewOAuth2GoogleHandler(svc google.AuthService, userSvc service.UserService) *OAuth2GoogleHandler {
	return &OAuth2GoogleHandler{
		googleAuthSvc:   svc,
		userSvc:         userSvc,
		key:             []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgB"),
		stateCookieName: "jwt-state",
	}
}

func (h *OAuth2GoogleHandler) RegisterRoutes(s *gin.Engine) {
	g := s.Group("/oauth2/google")
	g.GET("/login", h.Login)
	g.Any("callback", h.Callback)
}

func (h *OAuth2GoogleHandler) Login(ctx *gin.Context) {
	state, err := h.googleAuthSvc.NewState(32)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to generate state")
		return
	}
	err = h.setStateCookie(ctx, state)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to set cookie")
		return
	}
	ctx.Redirect(http.StatusFound, h.googleAuthSvc.AuthCodeURL(state))
}

func (h *OAuth2GoogleHandler) Callback(ctx *gin.Context) {
	if ctx.Request.Method != http.MethodGet {
		ctx.String(http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	err := h.verifyState(ctx)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid state")
		return
	}

	code := ctx.Query("code")
	if code == "" {
		ctx.String(http.StatusBadRequest, "Invalid request")
		return
	}

	tok, err := h.googleAuthSvc.ExchangeCode(ctx.Request.Context(), code)
	if err != nil {
		ctx.String(http.StatusBadRequest, "exchange failed: %v", err)
		return
	}

	httpClient := h.googleAuthSvc.Client(ctx.Request.Context(), tok)
	googleUser, err := h.googleAuthSvc.FetchUserInfo(ctx.Request.Context(), httpClient)
	if err != nil {
		ctx.String(http.StatusBadRequest, "userinfo failed: %v", err)
		return
	}

	user, err := h.userSvc.FindOrCreateByGoogle(ctx, googleUser)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to find or create user: %v", err)
		return
	}

	h.setJWTToken(ctx, user.ID)

	ctx.JSON(http.StatusOK, user)
}

func (h *OAuth2GoogleHandler) setStateCookie(ctx *gin.Context, state string) error {
	claims := stateClaims{
		State: state,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.key)
	if err != nil {
		return err
	}
	ctx.SetCookie(h.stateCookieName, tokenString, 0,
		"/oauth2/google/callback", "", false, true)
	return nil
}

func (h *OAuth2GoogleHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie(h.stateCookieName)
	if err != nil {
		return fmt.Errorf("cookie not found")
	}
	token, err := jwt.ParseWithClaims(ck, &stateClaims{}, func(token *jwt.Token) (interface{}, error) {
		return h.key, nil
	})
	if err != nil {
		return fmt.Errorf("invalid token: %v", err)
	}
	claims, ok := token.Claims.(*stateClaims)
	if !ok || claims.State != state {
		return fmt.Errorf("state not match")
	}
	return nil
}

type stateClaims struct {
	jwt.RegisteredClaims
	State string `json:"state"`
}

func (h *OAuth2GoogleHandler) Logout(ctx *gin.Context) {}
