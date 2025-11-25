package web

import (
	"fire/internal/service"
	"fire/internal/service/oauth2/google"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OAuth2GoogleHandler struct {
	googleAuthSvc google.AuthService
	userSvc       service.UserService
}

func NewOAuth2GoogleHandler(svc google.AuthService, userSvc service.UserService) *OAuth2GoogleHandler {
	return &OAuth2GoogleHandler{
		googleAuthSvc: svc,
		userSvc:       userSvc,
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
	ctx.Redirect(http.StatusFound, h.googleAuthSvc.AuthCodeURL(state))
}

func (h *OAuth2GoogleHandler) Callback(ctx *gin.Context) {
	if ctx.Request.Method != http.MethodGet {
		ctx.String(http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	state := ctx.Query("state")
	code := ctx.Query("code")
	if state == "" || code == "" {
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

	ctx.JSON(http.StatusOK, user)
}
