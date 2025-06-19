package controllerImpl

import (
	controller "github.com/berserkkv/trader/controller/interface"
	"github.com/berserkkv/trader/httpctx"
	service "github.com/berserkkv/trader/service/interface"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type OrderControllerImpl struct {
	s service.OrderService
}

func NewOrderControllerImpl(service service.OrderService) *OrderControllerImpl {
	return &OrderControllerImpl{s: service}
}

func (o *OrderControllerImpl) GetAll(c httpctx.Context) {
	orders := o.s.GetAll(nil)
	c.JSON(http.StatusOK, orders)
}

func (o *OrderControllerImpl) GetStatisticsByBotID(c httpctx.Context) {
	botIdParam := c.Query("botId")
	if botIdParam == "" {
		statsMap := o.s.GetAllStatistics()
		c.JSON(http.StatusOK, statsMap)
		return
	}

	botId, err := strconv.ParseInt(botIdParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "botId must be a valid integer"})
		return
	}

	s := o.s.GetStatisticsByBotID(botId)

	c.JSON(200, s)
}

func (o *OrderControllerImpl) Create(c httpctx.Context) {
	//var order model.Order
	//
	//if err := c.ShouldBindJSON(&order); err != nil {
	//	fmt.Println("Bind Error:", err)
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//
	//newOrder := service.CreateOrder(order)
	//c.JSON(http.StatusCreated, newOrder)
}
func (o *OrderControllerImpl) Update(c httpctx.Context) {
	//var order model.Order
	//if err := c.ShouldBind(&order); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//updated := repository.UpdateOrder(order)
	//c.JSON(http.StatusCreated, updated)
}
func (o *OrderControllerImpl) GetAllByBotID(c httpctx.Context) {
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

	orders := o.s.GetAllByBotID(botId)
	c.JSON(200, orders)
}

var _ controller.OrderController = (*OrderControllerImpl)(nil)
