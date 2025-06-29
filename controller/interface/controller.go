package controller

import (
	"github.com/berserkkv/trader/httpctx"
)

type BotController interface {
	GetAllBots(ctx httpctx.Context)
	GetBotByID(c httpctx.Context)
	CreateBot(c httpctx.Context)
	DeleteBot(c httpctx.Context)
	StartBot(c httpctx.Context)
	StopBot(c httpctx.Context)
	ClosePosition(c httpctx.Context)
	Update(c httpctx.Context)
}

type OrderController interface {
	GetAll(c httpctx.Context)
	GetStatisticsByBotID(c httpctx.Context)
	GetAllStatistics(c httpctx.Context)
	Create(c httpctx.Context)
	Update(c httpctx.Context)
	GetAllByBotID(c httpctx.Context)
}
