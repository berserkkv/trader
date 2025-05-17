package controller

import (
	"bytes"
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/repository"
	"github.com/berserkkv/trader/service"
	"github.com/berserkkv/trader/service/connector"
	logger "github.com/berserkkv/trader/tools/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

var log *slog.Logger

func printBody(c *gin.Context) {
	b, _ := io.ReadAll(c.Request.Body)

	fmt.Println(string(b))

	c.Request.Body = io.NopCloser(bytes.NewBuffer(b))
}

func Register() {
	log = logger.Get()
	r := gin.Default()

	r.Use(cors.Default())

	api := r.Group("/api/users")
	{
		api.GET("", GetUsers)
		api.GET("/:id", GetUser)
		api.POST("", CreateUser)
		api.PUT("", UpdateUser)
		api.DELETE("/:id", DeleteUser)
	}

	positions := r.Group("/api/positions")
	{
		positions.GET("", GetPositions)
		positions.GET("/:id", GetPositionById)
		positions.POST("", CreatePosition)
		positions.PUT("", UpdatePosition)
		positions.DELETE("/:id", DeletePositionById)
	}

	price := r.Group("/api/prices")
	{
		price.GET("", GetSolPrice)
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

func GetPositions(c *gin.Context) {
	positions := service.GetAllPosition()
	c.JSON(http.StatusOK, positions)
}

func GetPositionById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	position, error := repository.GetPositionById(uint(id))
	if error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusOK, position)

}

func CreatePosition(c *gin.Context) {
	var position model.Position

	if err := c.ShouldBindJSON(&position); err != nil {
		fmt.Println("Bind Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPosition := service.CreatePosition(position)
	c.JSON(http.StatusCreated, newPosition)
}

func UpdatePosition(c *gin.Context) {
	var position model.Position
	if err := c.ShouldBind(&position); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated := repository.UpdatePosition(position)
	c.JSON(http.StatusCreated, updated)
}

func DeletePositionById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	service.DeletePositionById(uint(id))
}

func GetUsers(c *gin.Context) {
	users := service.GetAll()
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUser := service.Create(user)
	c.JSON(http.StatusCreated, newUser)
}

func UpdateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated := service.Update(user)
	c.JSON(http.StatusOK, updated)
}

func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	service.Delete(uint(id))
	c.Status(http.StatusNoContent)
}
