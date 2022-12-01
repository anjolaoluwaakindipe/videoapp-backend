package services;

type HelloService struct {}

func (hs HelloService) SayHello () string {
	return "Hello"
}

func NewHelloService() HelloService {
	return HelloService{}
}