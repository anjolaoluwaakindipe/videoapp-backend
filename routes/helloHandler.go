package routes

import (
	"encoding/json"
	"net/http"

	"github.com/anjolaoluwaakindipe/videoapp/services"
)

type HelloHandler struct {
	helloService services.HelloServiceInterface
}

func (hh HelloHandler) GetHello() RouteHandler {
	return RouteHandler{
		Path:    "/hello",
		Methods: []string{"GET"},
		HandlerFunc: func(rw http.ResponseWriter, r *http.Request) {
			r.Header.Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			json.NewEncoder(rw).Encode(map[string]interface{}{
				"hello": hh.helloService.SayHello(),
			})
		},
	}
}

func NewHelloHandler(helloService services.HelloServiceInterface ) HelloHandler {
	return HelloHandler{helloService: helloService}
}
