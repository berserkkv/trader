package service

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/repository"
	"github.com/berserkkv/trader/service/connector"
	"time"
)

func GetAllOrders() []model.Order {
	return repository.GelAllOrders()
}

func CreateOrder(order model.Order) model.Order {

	now := time.Now()
	order.CreatedAt = now
	order.Symbol = "SOLUSDT"
	order.Side = "BUY"
	order.Type = "MARKET"
	order.Price = connector.GetPrice("SOLUSDT")
	order.Quantity = 100

	return repository.CreateOrder(order)
}

func UpdateOrder(order model.Order) model.Order {
	order.FilledAt = time.Now()
	order.Status = "Closed"
	return repository.UpdateOrder(order)
}
