package web

import (
	"bytes"
	"encoding/json"
	"fire/internal/domain"
	"fire/internal/service"
	svcmock "fire/internal/service/mocks"
	"fire/internal/web/jwt"
	"fire/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestArticleHandler_Publish(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.ArticleService
		reqBody  string
		wantCode int
		wantRes  Result
	}{
		{
			name: "publish_article",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmock.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "Test Title",
					Content: "Test Content",
					Author: domain.Author{
						ID: 123,
					},
				}).Return(int64(1), nil)
				return svc
			},
			reqBody: `
{
	"title": "Test Title",
	"content": "Test Content"
}`,
			wantCode: http.StatusOK,
			wantRes:  Result{Code: 0, Msg: "success", Data: float64(1)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			artSvc := tc.mock(ctrl)
			log := logger.NewNopLogger()
			handle := NewArticleHandler(artSvc, log)

			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("user", jwt.UserClaims{
					UserID: 123,
				})
			})
			handle.RegisterRoutes(server)
			req, err := http.NewRequest("POST", "/articles/publish", bytes.NewBufferString(tc.reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			var res Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
