package services

import "go.uber.org/fx"

var Module = fx.Module("services", fx.Provide(
	fx.Annotate(NewHelloService, fx.As(new(HelloServiceInterface))),
	fx.Annotate(NewRoomService, fx.As(new(RoomServiceInterface))),
	))
