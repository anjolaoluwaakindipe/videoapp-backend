package routes

import "github.com/anjolaoluwaakindipe/videoapp/utils/logger"

type RoomRoute struct {
	logger      logger.Logger
	roomHandler RoomHandler
}

// registers all room routes
func (rr RoomRoute) Routes() []RouteHandler {
	rr.logger.Infoln("Registering Room Route")

	// room routes slice tor return
	routes := []RouteHandler{
		rr.roomHandler.CreateRoom(),
		rr.roomHandler.JoinRoom(),
	}

	return routes
}

func NewRoomRoute(logger logger.Logger, roomHandler RoomHandler) RoomRoute {
	return RoomRoute{logger: logger, roomHandler: roomHandler}
}
