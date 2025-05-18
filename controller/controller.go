package controller

import (
	"bytes"
	"fmt"
	"github.com/berserkkv/trader/service/connector"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"net/http"
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

	price := r.Group("/api/prices")
	{
		price.GET("", GetSolPrice)
	}

	bots := r.Group("/api/bots")
	{
		bots.GET("", GetAllBots)
		bots.POST("", CreateBot)
		bots.PATCH("/stop/:id", StopBot)
		bots.PATCH("/start/:id", StartBot)
	}

	orders := r.Group("/api/orders")
	{
		orders.GET("", GetOrders)
		orders.POST("", CreateOrder)
		orders.PUT("", UpdateOrder)

	}

	r.Run(":8080")
}

func GetSolPrice(c *gin.Context) {
	c.JSON(http.StatusOK, connector.GetPrice("SOLUSDT"))
}
