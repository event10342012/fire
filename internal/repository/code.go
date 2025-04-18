package repository

import (
	"context"
	"fire/internal/repository/cache"
)

var ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany

type CodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(cache cache.CodeCache) CodeRepository {
	return CodeRepository{
		cache: cache,
	}
}

func (repo *CodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}

func (repo *CodeRepository) Get(ctx context.Context, biz, phone string) (string, error) {
	return repo.cache.Get(ctx, biz, phone)
}

func (repo *CodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, code)
}
