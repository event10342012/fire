package integration

import (
	"bytes"
	"encoding/json"
	"fire/internal/integration/startup"
	"fire/internal/repository/dao"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArticleHandler_Edit(t *testing.T) {
	db := startup.InitDB()
	server := startup.InitWebserver()
	defer func() {
		db.Exec("truncate table `articles`")
	}()

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
				err := db.Where("author_id = ?", 1).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.ID > 0)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Mtime > 0)
				assert.Equal(t, art.Title, "test title")
				assert.Equal(t, art.Content, "test content")
			},
			art: Article{
				Title:   "test title",
				Content: "test content",
			},
			wantCode: http.StatusOK,
			WantRes: Result[int64]{
				Code: 0,
				Msg:  "success",
				Data: 1,
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
			server.ServeHTTP(recoder, req)
			assert.Equal(t, tc.wantCode, recoder.Code)
			if recoder.Code != http.StatusOK {

			}
			var res Result[int64]
			err = json.Unmarshal(recoder.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, tc.WantRes, res)
		})
	}
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
