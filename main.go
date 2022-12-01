package main

import (
	"github.com/anjolaoluwaakindipe/videoapp/routes"
	"github.com/anjolaoluwaakindipe/videoapp/server"
	"github.com/anjolaoluwaakindipe/videoapp/services"
	"github.com/anjolaoluwaakindipe/videoapp/utils/config"
	"github.com/anjolaoluwaakindipe/videoapp/utils/logger"
	"go.uber.org/fx"
)

func main() {

	fx.New(
		server.Module, 
		config.Module, 
		logger.Module, 
		routes.RouteModule, 
		routes.HandlerModule,
		services.Module,
	).Run()
}
