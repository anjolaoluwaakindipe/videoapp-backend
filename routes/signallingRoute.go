package routes

import "github.com/anjolaoluwaakindipe/videoapp/utils/logger"

type SignallingRoute struct {
	logger logger.Logger
	signallingHandler SignallingHandler
}

func (sr SignallingRoute) Routes() []RouteHandler {
	sr.logger.Info("Registering Signalling Route")

	routes := []RouteHandler{

	}

	return routes
}

func NewSignallingRoute (logger logger.Logger, signallingHandler SignallingHandler) SignallingRoute{
	return SignallingRoute{logger: logger, signallingHandler: signallingHandler}
}