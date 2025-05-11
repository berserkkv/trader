package controller

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/service"
	logger "github.com/berserkkv/trader/tools/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

var log *slog.Logger

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello")
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

	r.Run(":8080")
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
