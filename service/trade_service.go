package service

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/repository"
)

func GetAllPosition() []model.Position {
	return repository.GetAllPositions()
}

func GetPositionById(id uint) (model.Position, error) {
	return repository.GetPositionById(id)
}

func CreatePosition(position model.Position) model.Position {
	return repository.CreatePosition(position)
}

func UpdatePosition(position model.Position) model.Position {
	return repository.UpdatePosition(position)
}

func DeletePositionById(id uint) {
	repository.DeletePositionById(id)
}

func GetAll() []model.User {
	return repository.GetAllUsers()
}

func GetByID(id uint) (model.User, error) {
	return repository.GetUserByID(id)
}

func Create(user model.User) model.User {
	return repository.CreateUser(user)
}

func Update(user model.User) model.User {
	return repository.UpdateUser(user)
}

func Delete(id uint) {
	repository.DeleteUser(id)
}
