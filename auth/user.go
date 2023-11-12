package auth

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email       string `json:"email" gorm:"unique"`
	Password    string `json:"password"`
	GivenName   string `json:"given_name"`
	FamilyName  string `json:"family_name"`
	Picture     string `json:"picture"`
	Locale      string `json:"locale"`
	GoogleId    string `json:"googleId"`
	IsSuperUser bool
	IsActive    bool
}

func (user *User) GetFullName() string {
	return user.GivenName + user.FamilyName
}

func (user *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

func (user *User) CheckPassword(hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(user.Password))
	return err == nil
}

func (user *User) Read(db *gorm.DB, id int) error {
	result := db.Take(user, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
