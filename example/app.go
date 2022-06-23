package main

import (
	"fmt"
	"strconv"
	"time"

	service "github.com/budimanlai/go-cli-service"
)

var (
	isBreak bool
	node    string
)

func StartFunc(context service.ServiceContext) {
	isBreak = false
	node = strconv.Itoa(context.Node)

	log("Service run as node: " + node)

	db := context.Database()
	result, e := db.Get("SELECT id, handphone FROM user WHERE handphone = ?", "62813813825254")
	if e != nil {
		log(e.Error())
	}

	if result == nil {
		log("user not found")
	} else {
		log("User ID: " + result.String("id") + ", Phone: " + result.String("handphone"))
	}

	count := 0
	for {
		log("Infinite loop: " + strconv.Itoa(count))
		time.Sleep(time.Second * 2)
		count++

		if isBreak {
			log("Service stoped")
			break
		}
	}
}

func StopFunc(context service.ServiceContext) {
	log("try to stop node")
	defer func() {
		isBreak = true
	}()
}

func log(msg string) {
	t := time.Now()

	s := "[Node: " + node + " " + t.Format("2006-01-02 15:04:05") + "] " + msg
	fmt.Println(s)
}
