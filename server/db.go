package server

import (
	"fire/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	dsn := "root@tcp(127.0.0.1:3306)/fire?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(model.User{}, model.AccountingTransaction{}, model.AssetPosition{}, model.Stock{})
	if err != nil {
		return
	}
}

func GetDB() *gorm.DB {
	return db
}
