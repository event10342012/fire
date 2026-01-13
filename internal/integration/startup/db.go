package startup

import (
	"fire/internal/repository/dao"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open("host=localhost user=leochen dbname=postgres search_path=fire port=5432 sslmode=disable TimeZone=Asia/Taipei"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
