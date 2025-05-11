package database

import (
	"github.com/berserkkv/trader/model"
	"gorm.io/gorm"

	"github.com/glebarez/sqlite"
)

var DB *gorm.DB

func ConnectDB() {
	db, err := gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.User{})
	DB = db
}
