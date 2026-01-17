package integration

import (
	"bytes"
	"encoding/json"
	"fire/internal/integration/startup"
	"fire/internal/repository/dao"
	"fire/internal/web/jwt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ArticleHandlerSuite struct {
	db     *gorm.DB
	server *gin.Engine
	suite.Suite
}

func (suite *ArticleHandlerSuite) SetupSuite() {
	suite.db = startup.InitDB()
	hdl := startup.InitArticleHandler()
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("user", jwt.UserClaims{
			UserID: 123,
		})
	})
	hdl.RegisterRoutes(server)
	suite.server = server
}

func (suite *ArticleHandlerSuite) TearDownSuite() {
	suite.db.Exec("truncate table fire.articles")
}

func (suite *ArticleHandlerSuite) TestEdit() {
	t := suite.T()

	testCases := []struct {
		name   string
		before func(*testing.T)
		after  func(*testing.T)
		art    Article

		wantCode int
		WantRes  Result[int64]
	}{
		{
			name: "create article success",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art dao.Article
				err := suite.db.Where("author_id = ?", 123).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.ID > 0)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Mtime > 0)
				assert.Equal(t, art.Title, "test title")
				assert.Equal(t, art.Content, "test content")
				assert.Equal(t, art.AuthorID, int64(123))
			},
			art: Article{
				Title:   "test title",
				Content: "test content",
			},
			wantCode: http.StatusOK,
			WantRes: Result[int64]{
				Code: 0,
				Msg:  "success",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest("POST", "/articles/edit", bytes.NewReader(reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			recoder := httptest.NewRecorder()
			suite.server.ServeHTTP(recoder, req)
			assert.Equal(t, tc.wantCode, recoder.Code)
			if recoder.Code != http.StatusOK {

			}
			var res Result[int64]
			err = json.Unmarshal(recoder.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, tc.WantRes.Msg, res.Msg)
			assert.Equal(t, tc.WantRes.Code, res.Code)
		})
	}
}

func TestArticleHandlerSuite(t *testing.T) {
	suite.Run(t, &ArticleHandlerSuite{})
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
