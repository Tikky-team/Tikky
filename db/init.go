package db

import (
	"Tikky/db/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDatabase() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/tikky?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("db connection failed: %v", err))
	}

	Db = db

	err = db.AutoMigrate(model.User{}, model.Video{}, model.Comment{}, model.Favorite{})
	if err != nil {
		panic(fmt.Errorf("db migrate failed: %v", err))
	}
}
