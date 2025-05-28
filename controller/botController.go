package controller

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/service"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

func GetAllBots(c *gin.Context) {
	fields := map[string]interface{}{}

	if val := c.Query("isNotActive"); val != "" {
		isNotActive, err := strconv.ParseBool(val)
		if err == nil {
			fields["is_not_active"] = isNotActive
		}
	}

	if strategy := c.Query("strategy"); strategy != "" {
		fields["strategy_name"] = strategy
	}

	if timeframe := c.Query("timeframe"); timeframe != "" {
		fields["timeframe"] = timeframe
	}

	if symbol := c.Query("symbol"); symbol != "" {
		fields["symbol"] = symbol
	}

	if val := c.Query("inPos"); val != "" {
		inPos, err := strconv.ParseBool(val)
		if err == nil {
			fields["in_pos"] = inPos
		}
	}

	bots := service.GetAllBots(fields)
	c.JSON(http.StatusOK, bots)
}

func GetBotById(c *gin.Context) {
	id, ok := extractId(c)
	if !ok {
		return
	}
	b := service.GetBotById(id)
	c.JSON(http.StatusOK, b)
}

func StopBot(c *gin.Context) {
	botId, ok := extractId(c)
	if !ok {
		return
	}
	err := service.StopBot(int64(botId))
	if err != nil {
		slog.Error("Error stopping bot", "id", botId, "error", err)
	}

	c.JSON(200, gin.H{"message": "Stoped bot", "id": botId})
}

func StartBot(c *gin.Context) {
	botId, ok := extractId(c)
	if !ok {
		return
	}

	err := service.StartBot(botId)
	if err != nil {
		slog.Error("Error starting bot", "id", botId, "error", err)
	}
	c.JSON(200, gin.H{"message": "Started bot", "id": botId})
}

func CreateBot(c *gin.Context) {
	var req model.CreateBotRequest

	if err := c.BindJSON(&req); err != nil {
		slog.Error("Invalid input for creating bot request", "error", err)
		c.JSON(400, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	bot, err := service.CreateBot(req.Symbol, req.Strategy, req.Timeframe, req.Capital, req.Leverage, req.TakeProfit, req.StopLoss)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not create bot", "details": err.Error()})
		return
	}
	c.JSON(200, bot)
}

func extractId(c *gin.Context) (int64, bool) {
	id := c.Param("id")
	botID, err := strconv.Atoi(id)
	if err != nil {
		slog.Error("error converting id to int", "error", err, id, "id", id)
		c.JSON(400, gin.H{"error": "Invalid bot ID"})
		return 0, false
	}
	return int64(botID), true
}
