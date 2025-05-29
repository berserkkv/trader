package controller

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/repository"
	"github.com/berserkkv/trader/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetOrders(c *gin.Context) {
	orders := service.GetAllOrders()
	c.JSON(http.StatusOK, orders)
}

func GetOrderStatistics(c *gin.Context) {
	botIdParam := c.Query("botId")
	if botIdParam == "" {
		c.JSON(400, gin.H{"error": "botId query parameter is required"})
		return
	}

	botId, err := strconv.ParseInt(botIdParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "botId must be a valid integer"})
		return
	}

	s := service.GetOrderStatistics(botId)

	//r := make([]model.Statistics, 0)
	//for i := 0; i < 50; i++ {
	//	stat := model.Statistics{
	//		Pnl:  rand.Float64()*200 - 100,                      // Random float between -100 and +100
	//		Time: time.Now().Add(-time.Duration(i) * time.Hour), // Each entry 1 hour apart in the past
	//	}
	//	r = append(r, stat)
	//}
	//sort.Slice(r, func(i, j int) bool {
	//	return r[i].Time.Before(r[j].Time)
	//})
	c.JSON(200, s)
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

func GetOrdersByBotId(c *gin.Context) {
	botIdParam := c.Query("botId")
	if botIdParam == "" {
		c.JSON(400, gin.H{"error": "botId query parameter is required"})
		return
	}

	botId, err := strconv.ParseInt(botIdParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "botId must be a valid integer"})
		return
	}

	orders := service.GetOrdersByBotId(botId)
	c.JSON(200, orders)
}
