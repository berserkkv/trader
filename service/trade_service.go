package service

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/repository"
)

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
