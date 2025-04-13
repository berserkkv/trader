package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Message struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func main() {
	router := gin.Default()

	// Serve static HTML and JS
	router.Static("/templates", "./templates")
	router.GET("/", func(c *gin.Context) {
		c.File("./templates/index.html")
	})

	// WebSocket handler to get real-time updates
	router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}
		defer conn.Close()

		for {
			// Get the symbol from the client
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}

			var symbol string
			if err := json.Unmarshal(msg, &symbol); err != nil {
				log.Println("Unmarshal error:", err)
				continue
			}

			// Connect to Binance WebSocket
			go fetchBinancePrice(symbol, conn)
		}
	})

	router.Run(":8080")
}

// Fetch price from Binance WebSocket
func fetchBinancePrice(symbol string, conn *websocket.Conn) {
	url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@trade", symbol)

	// Connect to Binance WebSocket
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("WebSocket connection error:", err)
		return
	}
	defer c.Close()

	for {
		// Read Binance WebSocket message
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			return
		}

		// Extract price from Binance message
		var data map[string]interface{}
		if err := json.Unmarshal(msg, &data); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		price := data["p"].(string)
		message := Message{Symbol: symbol, Price: price}

		// Send price update to the client via WebSocket
		if err := conn.WriteJSON(message); err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}
