package router

import (
	"bytes"
	"fmt"
	controller "github.com/berserkkv/trader/controller/interface"
	"github.com/berserkkv/trader/controller/pairBotController"
	"github.com/berserkkv/trader/controller/pairOrderController"
	"github.com/berserkkv/trader/httpctx/ctxImpl"
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

func Register(botController controller.BotController) {
	r := gin.Default()

	r.Use(cors.Default())

	bots := r.Group("/api/bots")
	{
		bots.GET("", func(c *gin.Context) {
			botController.GetAllBots(&ctxImpl.GinContext{C: c})
		})
		bots.POST("", func(c *gin.Context) {
			botController.CreateBot(&ctxImpl.GinContext{C: c})
		})
		bots.GET("/:id", func(c *gin.Context) {
			botController.GetBotByID(&ctxImpl.GinContext{C: c})
		})
		bots.PATCH("/:id/stop", func(c *gin.Context) {
			botController.StopBot(&ctxImpl.GinContext{C: c})
		})
		bots.PATCH("/:id/start", func(c *gin.Context) {
			botController.StartBot(&ctxImpl.GinContext{C: c})
		})
		bots.PATCH("/:id/close_position", func(c *gin.Context) {
			botController.ClosePosition(&ctxImpl.GinContext{C: c})
		})
		bots.DELETE("/:id", func(c *gin.Context) {
			botController.DeleteBot(&ctxImpl.GinContext{C: c})
		})
	}
	//orders := r.Group("/api/orders")
	//{
	//	orders.GET("", GetOrders)
	//	orders.POST("", CreateOrder)
	//	orders.PUT("", UpdateOrder)
	//	orders.GET("/by-bot", GetOrdersByBotId)
	//	orders.GET("/statistics", GetOrderStatistics)
	//
	//}

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
