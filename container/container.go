package container

import (
	"github.com/berserkkv/trader/bot/botFather"
	"github.com/berserkkv/trader/controller/controllerImpl"
	controller "github.com/berserkkv/trader/controller/interface"
	"github.com/berserkkv/trader/database"
	sqliteImpl "github.com/berserkkv/trader/repository/impl/sqlite"
	repository "github.com/berserkkv/trader/repository/interface"
	serviceImpl "github.com/berserkkv/trader/service/impl"
	service "github.com/berserkkv/trader/service/interface"
	"github.com/berserkkv/trader/tools/config"
	logger "github.com/berserkkv/trader/tools/log"
)

type Container struct {
	Config          *config.Config
	BotFather       *botFather.BotFather
	BotRepository   repository.BotRepository
	OrderRepository repository.OrderRepository
	BotService      service.BotService
	OrderService    service.OrderService
	BotController   controller.BotController
	OrderController controller.OrderController
}

func New() *Container {
	// Load config and initialize logger
	cnf := config.Load()
	logger.Init(cnf.Logger.Level, cnf.Env)

	// Initialize DB
	database.ConnectDB()
	db := database.DB

	// Repositories
	botRepo := sqliteImpl.NewBotRepository(db)
	orderRepo := sqliteImpl.NewOrderRepositoryImpl(db)

	// Instantiate core components
	bf := botFather.GetBotFather(botRepo, orderRepo)

	// Services
	botService := serviceImpl.NewBotService(botRepo, bf)
	orderService := serviceImpl.NewOrderServiceImpl(orderRepo)

	// Controllers
	botController := controllerImpl.NewBotController(botService)
	orderController := controllerImpl.NewOrderControllerImpl(orderService)

	return &Container{
		Config:          cnf,
		BotFather:       bf,
		BotRepository:   botRepo,
		OrderRepository: orderRepo,
		BotService:      botService,
		OrderService:    orderService,
		BotController:   botController,
		OrderController: orderController,
	}
}
