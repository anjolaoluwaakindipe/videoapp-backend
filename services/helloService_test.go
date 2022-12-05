package services_test

import (
	"testing"

	"github.com/anjolaoluwaakindipe/videoapp/services"
)

func TestSayHello(t *testing.T) {
	// test class
	helloService := services.NewHelloService()

	// test table 
	var testTable = []struct {
		name   string
		output string
	}{
		{name: "basic-test", output: "Hello"},
	}

	// test
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			result := helloService.SayHello()
			if test.output != result {
				t.Errorf(`Error in HelloService.SayHello(), expected "%v" got "%v"  `, result, test.output)
			}else{
				t.Logf("HelloService.SayHello() PASSED")
			}
		})
	}
}
