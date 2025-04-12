package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func handleTrade(c *gin.Context) {
	c.String(200, "Welcome to your trading server!")
}

func main() {
	// Initialize the Gin router
	r := gin.Default()

	// Define a route for the trade handler
	r.GET("/", handleTrade)

	// Start the server on port 8080
	fmt.Println("Server started at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
