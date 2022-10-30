package main

import (
	services "github.com/budimanlai/go-cli-service"
)

func main() {
	srv := services.NewService("config/main.conf")
	srv.AppName = "Example command line service"
	srv.Version = "0.0.1"
	srv.StartHandler(StartService)
	srv.StopHandler(StopService)
	e := srv.Start()
	if e != nil {
		panic(e)
	}
}
