package main

import (
	"time"

	services "github.com/budimanlai/go-cli-service"
)

func StartService(ctx *services.Service) {
	ctx.Log("Run with args:", ctx.Args.GetRawArgs())

	result, e := ctx.Db.Get("select version() as versi")
	if e != nil {
		ctx.Log("Error:", e.Error())
	}

	ctx.Log("Versi DB:", result.String("versi"))

	ctx.Ping.Start()

	for {
		ctx.Ping.Update()
		ctx.Log("Sleep...")
		time.Sleep(2 * time.Second)

		if ctx.IsStopped {
			ctx.Ping.Stop()
			ctx.Log("Exit from loop StartService")
			break
		}
	}
}

func StopService(ctx *services.Service) {
	ctx.Log("Stop Service")
	ctx.IsStopped = true
}
