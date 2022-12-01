package routes

import (
	"encoding/json"
	"net/http"

	"go.uber.org/fx"
)

type Route interface {
	Routes() []RouteHandler
}

// What every method of a handler should return
// and also what every route's "Routes()" method should return
type RouteHandler struct {
	Path    string
	Methods []string
	http.HandlerFunc
}

// Wraps any struct that has an interface of "Route" with fx Annotations
// that allow it to be used by the server
func RegisterRoute(r interface{}) interface{} {
	return fx.Annotate(r, fx.As(new(Route)), fx.ResultTags(`group:"routes"`))
}

var RouteModule = fx.Module("routes", fx.Provide(
	RegisterRoute(NewHelloRoute),
	RegisterRoute(NewRoomRoute),
	RegisterRoute(NewSignallingRoute),
))

var HandlerModule = fx.Module("handlers", fx.Provide(
	NewHelloHandler,
	NewRoomHandler,
	NewSignallingHandler,
))

func res(rw http.ResponseWriter, body interface{}, code int) {
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(code)
	json.NewEncoder(rw).Encode(body)
}
