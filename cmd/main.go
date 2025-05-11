package main

import (
	"github.com/berserkkv/trader/controller"
	"github.com/berserkkv/trader/tools/config"
	logger "github.com/berserkkv/trader/tools/log"
)

func main() {
	// init
	cnf := config.Load()
	logger.Init(cnf.Logger.Level, cnf.Env)

	log := logger.Get()

	log.Info("Server started on port: 8080")

	controller.Register()

}
