package server

import "gorm.io/gorm"

type Model interface {
	Read(db *gorm.DB, id int) error
	Create(db *gorm.DB) error
	Update(db *gorm.DB) error
	Delete(db *gorm.DB) error
}
