package repository

import (
	"github.com/berserkkv/trader/database"
	"github.com/berserkkv/trader/model"
)

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
