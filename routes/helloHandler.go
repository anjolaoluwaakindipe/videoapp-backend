package routes

import (
	"encoding/json"
	"net/http"

	"github.com/anjolaoluwaakindipe/videoapp/services"
)

type HelloHandler struct {
	helloService services.HelloServiceInterface
}

func (hh HelloHandler) getHello() RouteHandler {
	return RouteHandler{
		Path:    "/hello",
		Methods: []string{"GET"},
		HandlerFunc: func(rw http.ResponseWriter, r *http.Request) {
			r.Header.Add("Content-Type", "application/json")
			json.NewEncoder(rw).Encode(map[string]interface{}{
				"hello": hh.helloService.SayHello(),
			})
		},
	}
}

func NewHelloHandler () HelloHandler {
	return HelloHandler{}
}
