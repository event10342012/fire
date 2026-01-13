package middleware

import (
	"bytes"
	"context"
	"io"

	"github.com/gin-gonic/gin"
)

type LogMiddlewareBuilder struct {
	logFun        func(ctx context.Context, l AccessLog)
	allowReqBody  bool
	allowRespBody bool
}

func (l *LogMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if len(path) > 1024 {
			path = path[:1024]
		}
		method := ctx.Request.Method
		al := AccessLog{
			Path:   path,
			Method: method,
		}

		if l.allowReqBody {
			body, _ := ctx.GetRawData()
			if len(body) > 2048 {
				al.ReqBody = string(body[:2048])
			} else {
				al.ReqBody = string(body)
			}
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
	}
}

type AccessLog struct {
	Path     string `json:"path"`
	Method   string `json:"method"`
	ReqBody  string `json:"req_body"`
	RespBody string `json:"resp_body"`
}
