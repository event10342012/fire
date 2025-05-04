package service

import (
	"context"
	"fire/internal/domain"
	"fire/internal/repository"
	repomock "fire/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordEncrypt(t *testing.T) {
	password := []byte("123@qaz")
	encrypted, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)
	println(string(encrypted))
	err = bcrypt.CompareHashAndPassword(encrypted, []byte("123@qaz"))
	assert.NoError(t, err)
}

func TestUserService_Login(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		ctx      context.Context
		email    string
		password string
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "Login Success",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepo := repomock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().FindByEmail(
					gomock.Any(),
					"test@gmail.com",
				).Return(domain.User{
					ID:       1,
					Email:    "test@gmail.com",
					Password: "$2a$10$hGh9QRuIGK7AZmGafLIEme/roPnttYeZ6i9D7xhILKHuUFhSySws.",
				}, nil)
				return userRepo
			},
			email:    "test@gmail.com",
			password: "123@qaz",
			wantUser: domain.User{
				ID:       1,
				Email:    "test@gmail.com",
				Password: "$2a$10$hGh9QRuIGK7AZmGafLIEme/roPnttYeZ6i9D7xhILKHuUFhSySws.",
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := tc.mock(ctrl)
			userSvc := NewUserService(repo)
			user, err := userSvc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantUser, user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
