package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email       string    `json:"email" gorm:"unique"`
	Password    string    `json:"password"`
	GivenName   string    `json:"given_name"`
	FamilyName  string    `json:"family_name"`
	Birthdate   time.Time `json:"birthdate"`
	Picture     string    `json:"picture"`
	Locale      string    `json:"locale"`
	GoogleId    string    `json:"googleId"`
	IsSuperUser bool
	IsActive    bool
}

func (user *User) GetFullName() string {
	return user.GivenName + user.FamilyName
}

func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (user *User) IsAuthenticated() bool {
	return user.IsActive
}
