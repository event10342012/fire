package web

import (
	"fire/internal/service/oauth2/google"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type OAuth2GoogleHandler struct {
	googleAuthSvc google.AuthService
}

func NewOAuth2GoogleHandler(svc google.AuthService) *OAuth2GoogleHandler {
	return &OAuth2GoogleHandler{
		googleAuthSvc: svc,
	}
}

func (h *OAuth2GoogleHandler) RegisterRoutes(s *gin.Engine) {
	g := s.Group("/oauth2/google")
	g.GET("/login", h.Login)
	g.Any("callback", h.Callback)
}

func (h *OAuth2GoogleHandler) Login(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	if sess.Get("user") != nil {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	state, err := h.googleAuthSvc.NewState(32)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to generate state")
		return
	}

	sess.Set(google.SessionKeyState, state)
	err = sess.Save()
	if err != nil {
		return
	}

	ctx.Redirect(http.StatusFound, h.googleAuthSvc.AuthCodeURL(state))
}

func (h *OAuth2GoogleHandler) Callback(ctx *gin.Context) {
	if ctx.Request.Method != http.MethodGet {
		ctx.String(http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	sess := sessions.Default(ctx)
	stateExpected, _ := sess.Get(google.SessionKeyState).(string)
	state := ctx.Query("state")
	code := ctx.Query("code")
	if state == "" || code == "" || stateExpected != state {
		ctx.String(http.StatusBadRequest, "Invalid request")
		return
	}

	tok, err := h.googleAuthSvc.ExchangeCode(ctx.Request.Context(), code)
	if err != nil {
		ctx.String(http.StatusBadRequest, "exchange failed: %v", err)
		return
	}

	httpClient := h.googleAuthSvc.Client(ctx.Request.Context(), tok)
	user, err := h.googleAuthSvc.FetchUserInfo(ctx.Request.Context(), httpClient)
	if err != nil {
		ctx.String(http.StatusBadRequest, "userinfo failed: %v", err)
		return
	}

	// Persist in session (store minimal info; tokens are sensitive)
	sess.Set(google.SessionKeyAccessTok, tok.AccessToken)
	if tok.RefreshToken != "" {
		sess.Set(google.SessionKeyRefreshTok, tok.RefreshToken)
	}
	if idt, _ := tok.Extra("id_token").(string); idt != "" {
		sess.Set(google.SessionKeyIDToken, idt)
	}
	sess.Set(google.SessionKeyUser, user)
	_ = sess.Save()

	ctx.Redirect(http.StatusFound, "/me")
}
