package service

import (
	"context"
	"fire/internal/domain"
	"fire/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
}

type articleService struct {
	readerRepo repository.ArticleReaderRepository
	authorRepo repository.ArticleAuthorRepository
}

func NewArticleService(readerRepo repository.ArticleReaderRepository, authorRepo repository.ArticleAuthorRepository) ArticleService {
	return &articleService{readerRepo: readerRepo, authorRepo: authorRepo}
}

func (s *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	if art.ID > 0 {
		err := s.readerRepo.Create(ctx, art)
		return art.ID, err
	}
	return 1, nil
}

func (s *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	id, err := s.authorRepo.Create(ctx, art)
	if err != nil {
		return 0, err
	}
	art.ID = id
	err = s.readerRepo.Create(ctx, art)
	return id, err
}
