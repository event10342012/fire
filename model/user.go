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

func (user *User) ReadByGoogleID(db *gorm.DB, googleID string) error {
	result := db.Take(user, "google_id = ?", googleID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (user *User) ReadByEmail(db *gorm.DB, email string) error {
	result := db.Take(user, "email = ?", email)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (user *User) Read(db *gorm.DB, id int) error {
	result := db.Take(user, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (user *User) Create(db *gorm.DB) error {
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (user *User) Update(db *gorm.DB) error {
	result := db.Updates(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (user *User) Delete(db *gorm.DB) error {
	result := db.Delete(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
