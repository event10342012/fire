package repository

import (
	"context"
	"fire/internal/domain"
)

type ArticleReaderRepository interface {
	Create(ctx context.Context, article domain.Article) error
}
