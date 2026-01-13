package repository

import (
	"context"
	"fire/internal/domain"
	"fire/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
}

type CacheArticleRepository struct {
	dao dao.ArticleDAO
}

func NewArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &CacheArticleRepository{dao: dao}
}

func (a *CacheArticleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	return a.dao.Insert(ctx, a.toEntity(article))
}

func (a *CacheArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorID: art.Author.ID,
	}
}
