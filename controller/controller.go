package controller

import (
	"bytes"
	"fmt"
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
	}

	orders := r.Group("/api/orders")
	{
		orders.GET("", GetOrders)
		orders.POST("", CreateOrder)
		orders.PUT("", UpdateOrder)
		orders.GET("/by-bot", GetOrdersByBotId)

	}

	r.Run(":8080")
}
