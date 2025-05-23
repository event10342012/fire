package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode string

	ErrCodeSentTooMany   = errors.New("sent too many times")
	ErrCodeVerifyTooMany = errors.New("verify too many times")
)

type CodeCache interface {
	Get(ctx context.Context, biz, phone string) (string, error)
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type RedisCodeCache struct {
	cmd redis.Cmdable
}

func NewCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		cmd: cmd,
	}
}

func (c *RedisCodeCache) Get(ctx context.Context, biz, phone string) (string, error) {
	key := c.Key(biz, phone)
	return c.cmd.Get(ctx, key).Result()
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.cmd.Eval(ctx, luaSetCode, []string{c.Key(biz, phone)}, code).Int()

	if err != nil {
		return err
	}

	switch res {
	case -2:
		return errors.New("code exists but dose not have expired")
	case -1:
		return ErrCodeSentTooMany
	default:
		return nil
	}
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	res, err := c.cmd.Eval(ctx, luaVerifyCode, []string{c.Key(biz, phone)}, code).Int()
	if err != nil {
		return false, err
	}

	switch res {
	case -1:
		return false, nil
	case -2:
		return false, ErrCodeVerifyTooMany
	default:
		return true, nil
	}
}

func (c *RedisCodeCache) Key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
