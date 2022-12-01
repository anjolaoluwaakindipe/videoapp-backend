package routes

import "github.com/anjolaoluwaakindipe/videoapp/utils/logger"

type HelloRoute struct {
	helloHandler HelloHandler
	logger logger.Logger
}

func (hr HelloRoute) Routes() []RouteHandler {
	hr.logger.Info("Registering Hello Route")
	routes := []RouteHandler{
		hr.helloHandler.getHello(),
	}
	return routes
}

func NewHelloRoute(helloHandler HelloHandler, logger logger.Logger) HelloRoute {
	return HelloRoute{helloHandler: helloHandler, logger: logger};
}
