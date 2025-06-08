package database

import (
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/pairBot"
	"gorm.io/gorm"

	"github.com/glebarez/sqlite"
)

var DB *gorm.DB

func ConnectDB() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.Migrator().DropTable(&model.Order{})
	db.Migrator().DropTable(&bot.Bot{})
	db.Migrator().DropTable(&model.PairOrder{})
	db.Migrator().DropTable(&pairBot.PairBot{})

	db.AutoMigrate(&model.Order{})
	db.AutoMigrate(&bot.Bot{})
	db.AutoMigrate(&model.PairOrder{})
	db.AutoMigrate(&pairBot.PairBot{})
	DB = db
}
