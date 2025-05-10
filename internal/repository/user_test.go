package repository

import (
	"context"
	"fire/internal/domain"
	"fire/internal/repository/cache"
	cachemock "fire/internal/repository/cache/mocks"
	"fire/internal/repository/dao"
	daomock "fire/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCacheUserRepository_FindByID(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
		ctx      context.Context
		id       int64
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "Find by id success hit cache",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				userDAO := daomock.NewMockUserDAO(ctrl)
				userCache := cachemock.NewMockUserCache(ctrl)
				userCache.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{
					ID: 1,
				}, nil)
				return userDAO, userCache
			},
			ctx:      context.Background(),
			id:       1,
			wantUser: domain.User{ID: 1},
			wantErr:  nil,
		},
		{
			name: "Find by id success but does not hit cache",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				userDAO := daomock.NewMockUserDAO(ctrl)
				userCache := cachemock.NewMockUserCache(ctrl)
				userCache.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{}, cache.ErrKeyNotExist)
				userCache.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil)

				userDAO.EXPECT().FindByID(gomock.Any(), int64(1)).Return(dao.User{
					ID:       1,
					Birthday: 100,
				}, nil)
				return userDAO, userCache
			},
			ctx:      context.Background(),
			id:       1,
			wantUser: domain.User{ID: 1, Birthday: time.UnixMilli(100)},
			wantErr:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userDao, userCache := tc.mock(ctrl)
			repo := NewUserRepository(userDao, userCache)
			user, err := repo.FindByID(tc.ctx, tc.id)
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, user, tc.wantUser)
		})
	}
}
