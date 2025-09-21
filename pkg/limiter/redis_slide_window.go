package limiter

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisSlideWindow struct {
	cmd redis.Cmdable
}

func (r RedisSlideWindow) Limit(ctx context.Context, key string) {
	//TODO implement me
	panic("implement me")
}
