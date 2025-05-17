package repository

import (
	"github.com/berserkkv/trader/database"
	"github.com/berserkkv/trader/model"
)

func GetAllPositions() []model.Position {
	var positions []model.Position
	database.DB.Find(&positions)
	return positions
}

func GetPositionById(id uint) (model.Position, error) {
	var position model.Position
	err := database.DB.First(&position, id).Error
	return position, err
}

func CreatePosition(position model.Position) model.Position {
	database.DB.Create(&position)
	return position
}

func UpdatePosition(position model.Position) model.Position {
	database.DB.Save(&position)
	return position
}

func DeletePositionById(id uint) {
	database.DB.Delete(&model.Position{}, id)
}

func GetAllUsers() []model.User {
	var users []model.User
	database.DB.Find(&users)
	return users
}

func GetUserByID(id uint) (model.User, error) {
	var user model.User
	err := database.DB.First(&user, id).Error
	return user, err
}

func CreateUser(user model.User) model.User {
	database.DB.Create(&user)
	return user
}

func UpdateUser(user model.User) model.User {
	database.DB.Save(&user)
	return user
}

func DeleteUser(id uint) {
	database.DB.Delete(&model.User{}, id)
}
