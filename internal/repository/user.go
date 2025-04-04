package repository

import (
	"context"
	"fire/internal/domain"
	"fire/internal/repository/dao"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), err
}

func (repo *UserRepository) FindByID(ctx context.Context, id int) (domain.User, error) {
	u, err := repo.dao.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	du := repo.toDomain(u)
	return du, nil
}

func (repo *UserRepository) UpdateNonZeroFields(ctx context.Context, user domain.User) error {
	u := dao.User{
		ID:       user.ID,
		Nickname: user.Nickname,
		Birthday: user.Birthday.UnixMilli(),
		AboutMe:  user.AboutMe,
	}
	return repo.dao.Update(ctx, u)
}

func (repo *UserRepository) toDomain(user dao.User) domain.User {
	return domain.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
	}
}
