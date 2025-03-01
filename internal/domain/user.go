package domain

import "time"

type User struct {
	ID          int64
	Email       string
	Password    string
	Phone       string
	GivenName   string
	FamilyName  string
	Nickname    string
	Birthday    time.Time
	AboutMe     string
	Picture     string
	Locale      string
	GoogleId    string
	IsSuperUser bool
	IsActive    bool
	CreatedAt   time.Time
}
