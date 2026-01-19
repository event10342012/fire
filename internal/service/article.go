package service

import (
	"context"
	"fire/internal/domain"
	"fire/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{repo: repo}
}

func (s *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	if art.ID > 0 {
		err := s.repo.Update(ctx, art)
		return art.ID, err
	}
	return s.repo.Create(ctx, art)
}
