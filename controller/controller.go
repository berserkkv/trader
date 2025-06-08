package controller

import (
	"bytes"
	"fmt"
	"github.com/berserkkv/trader/controller/pairBotController"
	"github.com/berserkkv/trader/controller/pairOrderController"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
)

var log *slog.Logger

func printBody(c *gin.Context) {
	b, _ := io.ReadAll(c.Request.Body)

	fmt.Println(string(b))

	c.Request.Body = io.NopCloser(bytes.NewBuffer(b))
}

func Register() {
	r := gin.Default()

	r.Use(cors.Default())

	bots := r.Group("/api/bots")
	{
		bots.GET("", GetAllBots)
		bots.POST("", CreateBot)
		bots.GET("/:id", GetBotById)
		bots.PATCH("/:id/stop", StopBot)
		bots.PATCH("/:id/start", StartBot)
		bots.PATCH("/:id/close_position", ClosePosition)
	}

	orders := r.Group("/api/orders")
	{
		orders.GET("", GetOrders)
		orders.POST("", CreateOrder)
		orders.PUT("", UpdateOrder)
		orders.GET("/by-bot", GetOrdersByBotId)
		orders.GET("/statistics", GetOrderStatistics)

	}

	pairBots := r.Group("/api/pair_bots")
	{
		pairBots.GET("", pairBotController.GetAllBots)
		pairBots.POST("", pairBotController.CreateBot)
		pairBots.GET("/:id", pairBotController.GetBotById)
		pairBots.PATCH("/:id/stop", pairBotController.StopBot)
		pairBots.PATCH("/:id/start", pairBotController.StartBot)
		pairBots.PATCH("/:id/close_position", pairBotController.ClosePosition)
	}

	pairOrders := r.Group("/api/pair_orders")
	{
		pairOrders.GET("", pairOrderController.GetOrders)
		pairOrders.POST("", pairOrderController.CreateOrder)
		pairOrders.PUT("", pairOrderController.UpdateOrder)
		pairOrders.GET("/by-bot", pairOrderController.GetOrdersByBotId)
		pairOrders.GET("/statistics", pairOrderController.GetOrderStatistics)

	}

	r.Run(":8080")
}
