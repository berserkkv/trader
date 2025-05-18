package service

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/repository"
)

func GetAllOrders() []model.Order {
	return repository.GelAllOrders()
}

func CreateOrder(order model.Order) model.Order {
	return repository.CreateOrder(order)
}

func UpdateOrder(order model.Order) model.Order {
	return repository.UpdateOrder(order)
}
