package main

import (
	service "github.com/budimanlai/go-service"
)

func main() {

	services := service.NewService()
	services.SetVersion("1.0.0 build 1234")
	services.Start(StartFunc)
	services.Stop(StopFunc)
	services.Run()
}
