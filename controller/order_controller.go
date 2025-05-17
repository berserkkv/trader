package controller

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/repository"
	"github.com/berserkkv/trader/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOrders(c *gin.Context) {
	orders := service.GetAllOrders()
	c.JSON(http.StatusOK, orders)
}

func CreateOrder(c *gin.Context) {
	var order model.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		fmt.Println("Bind Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newOrder := service.CreateOrder(order)
	c.JSON(http.StatusCreated, newOrder)
}

func UpdateOrder(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBind(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated := repository.UpdateOrder(order)
	c.JSON(http.StatusCreated, updated)
}
