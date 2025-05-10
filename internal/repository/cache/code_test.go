package cache

import (
	"context"
	"fire/internal/repository/cache/redismocks"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func keyFunc(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		ctx     context.Context
		biz     string
		phone   string
		code    string
		wantErr error
	}{
		{
			name: "set success",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(0))
				res.EXPECT().Eval(
					gomock.Any(),
					luaSetCode,
					[]string{keyFunc("test", "1234567890")},
					[]any{"123456"},
				).Return(cmd)
				return res
			},
			ctx:   context.Background(),
			biz:   "test",
			phone: "1234567890",
			code:  "123456",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cmd := tc.mock(ctrl)
			cache := NewCodeCache(cmd)
			err := cache.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, err, tc.wantErr)
		})
	}
}
