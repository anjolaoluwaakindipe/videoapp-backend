package routes_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anjolaoluwaakindipe/videoapp/routes"
	"github.com/anjolaoluwaakindipe/videoapp/services"
)

func TestHelloHandler_getHello(t *testing.T) {
	helloService := services.NewHelloService()
	helloHandler := routes.NewHelloHandler(helloService)

	getHelloRouteHandler := helloHandler.GetHello()

	// request set up
	req, err := http.NewRequest(getHelloRouteHandler.Methods[0], getHelloRouteHandler.Path, nil)

	if err != nil {
		t.Fatal(err)
	}

	// response recorder
	rr := httptest.NewRecorder()

	// handler stup
	handler := http.HandlerFunc(getHelloRouteHandler.HandlerFunc)

	// server set up
	handler.ServeHTTP(rr, req)

	// check status code
	if statusCode := rr.Code; statusCode != http.StatusOK {
		t.Errorf(`handler returned wrong status code. Expected "%v" got "%v"`, http.StatusOK, rr.Code)
	}

	expected, err := json.Marshal(map[string]string{
		"hello": "Hello",
	})

	if err != nil {
		t.Error(`An error occured while Marshalling expected json`)
	}

	if strings.EqualFold(string(expected), rr.Body.String()) {
		t.Errorf(`HelloHandler.GetHello handler failed, expected "%v" but got "%v"`, string(expected), rr.Body.String())
	} else {
		t.Log(`HelloHandler.GetHello PASSED`)
	}
}
