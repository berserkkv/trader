package controller

import (
	"github.com/berserkkv/trader/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllStatistics(ctx *gin.Context) {
	stats, err := service.GetAllStatistics()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, stats)
}

func GetStatisticsBySymbol(ctx *gin.Context) {
	symbol := ctx.Param("symbol")
	stats, err := service.GetStatisticsBySymbol(symbol)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, stats)
}

func GetStatisticsByTimeframe(ctx *gin.Context) {
	timeframe := ctx.Param("timeframe")
	stats, err := service.GetStatisticsByTimeframe(timeframe)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, stats)
}

func GetStatisticsByStrategy(ctx *gin.Context) {
	strategy := ctx.Param("strategy")
	stats, err := service.GetStatisticsByStrategy(strategy)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, stats)
}
