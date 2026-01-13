package dao

import (
	"context"

	"gorm.io/gorm"
)

type Article struct {
	ID       int64  `gorm:"primaryKey,autoIncrement"`
	Title    string `gorm:"type=varchar(4096)"`
	Content  string `gorm:"type=BLOB"`
	AuthorID int64  `gorm:"index"`
	Ctime    int64
	Mtime    int64
}

type ArticleDAO interface {
	Insert(ctx context.Context, article Article) (int64, error)
	//FindByID(ctx context.Context, id int64) (Article, error)
	//UpdateByID(ctx context.Context, article Article) error
	//DeleteByID(ctx context.Context, id int64) error
	//Count(ctx context.Context) (int64, error)
}

type ArticleGormDAO struct {
	db *gorm.DB
}

func NewArticleDAO(db *gorm.DB) ArticleDAO {
	return &ArticleGormDAO{
		db: db,
	}
}

func (dao *ArticleGormDAO) Insert(ctx context.Context, article Article) (int64, error) {
	result := dao.db.WithContext(ctx).Create(&article)
	if result.Error != nil {
		return 0, result.Error
	}
	return article.ID, nil
}
