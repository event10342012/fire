package web

import (
	"fire/internal/domain"
	"fire/internal/service"
	"fire/internal/web/jwt"
	"fire/pkg/logger"

	"net/http"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	svc service.ArticleService
	log logger.Logger
}

func NewArticleHandler(svc service.ArticleService, log logger.Logger) *ArticleHandler {
	return &ArticleHandler{svc: svc, log: log}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", h.Edit)
}

func (h *ArticleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		return
	}

	uc := ctx.MustGet("user").(jwt.UserClaims)
	id, err := h.svc.Save(ctx, domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			ID: uc.UserID,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 5,
			Msg:  "System error",
		})
		h.log.Error("save article error", logger.Error(err), logger.Int64("user_id", uc.UserID))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "success",
		Data: id,
	})
}
