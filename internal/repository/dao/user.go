package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateUser  = errors.New("user already exists")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, user User) error
	UpdateByID(ctx context.Context, user User) error
	Update(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, id int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
}

type GormUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GormUserDAO{
		db: db,
	}
}

func (dao *GormUserDAO) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.CreatedAt = now
	user.UpdatedAt = now
	err := dao.db.WithContext(ctx).Create(&user).Error
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		if me.Number == 1062 {
			return ErrDuplicateUser
		}
	}
	return err
}

func (dao *GormUserDAO) UpdateByID(ctx context.Context, user User) error {
	return dao.db.WithContext(ctx).Model(&user).Updates(User{
		Birthday:  user.Birthday,
		Nickname:  user.Nickname,
		AboutMe:   user.AboutMe,
		UpdatedAt: time.Now().UnixMilli(),
	}).Error
}

func (dao *GormUserDAO) Update(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.UpdatedAt = now
	return dao.db.WithContext(ctx).Save(user).Error
}

func (dao *GormUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (dao *GormUserDAO) FindByID(ctx context.Context, id int64) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return user, err
}

func (dao *GormUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	return user, err
}

type User struct {
	ID          int64          `gorm:"primaryKey autoIncrement"`
	Email       sql.NullString `gorm:"unique"`
	Password    string
	Phone       sql.NullString `gorm:"unique"`
	GivenName   string         `gorm:"type=varchar(128)"`
	FamilyName  string
	Nickname    string
	Birthday    int64
	AboutMe     string `gorm:"type=varchar(4096)"`
	Picture     string
	Locale      string
	GoogleId    string
	IsSuperUser bool
	IsActive    bool
	CreatedAt   int64
	UpdatedAt   int64
	DeletedAt   int64
}
