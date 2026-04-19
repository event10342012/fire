package service

import (
	"context"
	"fire/internal/domain"
	"fire/internal/repository"
	repomocks "fire/internal/repository/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_articleService_Publish(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleReaderRepository)
		art     domain.Article
		wantId  int64
		wantErr error
	}{
		{
			name: "Publish article Success",
			mock: func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleReaderRepository) {
				readerRepo := repomocks.NewMockArticleReaderRepository(ctrl)
				authorRepo := repomocks.NewMockArticleAuthorRepository(ctrl)
				authorRepo.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "Test Title",
					Content: "Test Content",
					Author: domain.Author{
						ID: 123,
					},
				}).Return(int64(1), nil).Times(1)
				readerRepo.EXPECT().Create(gomock.Any(), domain.Article{
					ID:      1,
					Title:   "Test Title",
					Content: "Test Content",
					Author: domain.Author{
						ID: 123,
					},
				}).Times(1)
				return authorRepo, readerRepo
			},
			art: domain.Article{
				Title:   "Test Title",
				Content: "Test Content",
				Author: domain.Author{
					ID: 123,
				},
			},
			wantId: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			authorRepo, readerRepo := tc.mock(ctrl)
			service := NewArticleService(readerRepo, authorRepo)
			id, err := service.Publish(context.Background(), tc.art)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantId, id)
		})
	}
}
