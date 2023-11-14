package model

import "gorm.io/gorm"

type BaseModel interface {
	Read(db *gorm.DB, id int) error
	Create(db *gorm.DB) error
	Update(db *gorm.DB) error
	Delete(db *gorm.DB) error
}
